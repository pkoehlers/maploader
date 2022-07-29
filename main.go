package main

import (
	"fmt"
	"log"
	"maploader/config"
	"maploader/robot"
	"maploader/tar"
	"maploader/util"
	"path/filepath"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const stateTopic string = "valetudo/maploader/status"
const currentMapTopic string = "valetudo/maploader/map"
const commandTopic string = currentMapTopic + "/set"

var maploaderDir string
var currentMap string

var roationKeepMaps int

var messageStateTopicHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	if string(msg.Payload()) != currentMap {
		log.Printf("Loaded current map from status topic")
		currentMap = string(msg.Payload())
	}
}

var messageCommandTopicHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	changeMap(client, string(msg.Payload()))
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")

	subCommandTopic(client)
	subCurrentMapTopic(client)
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}

func main() {
	currentMap = config.Getenv("DEFAULT_MAP_NAME", "main")
	maploaderDir = config.Getenv("MAPLOADER_DIR", "/data/maploader")
	roationKeepMaps = config.RotationKeepMaps()

	util.InitLogging(maploaderDir + "/log")
	config.InitConfig(config.Getenv("VALETUDO_CONFIG_PATH", "/data/valetudo_config.json"))

	robot.DetectRobot()

	var broker = config.MqttHost()
	var port = config.MqttPort()
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("maploader")
	opts.SetUsername(config.MqttUsername())
	opts.SetPassword(config.MqttPassword())
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	opts.WillTopic = stateTopic
	opts.WillEnabled = true
	opts.WillPayload = []byte("offline")
	opts.WillRetained = true

	client := mqtt.NewClient(opts)
	retry := time.NewTicker(5 * time.Second)
RetryLoop:
	for {
		select {
		case <-retry.C:
			if token := client.Connect(); token.Wait() && token.Error() != nil {
				//handle error
			} else {
				retry.Stop()
				break RetryLoop
			}
		}
	}
	publishState(client, "idle")
	for {
		time.Sleep(time.Second)
	}
}

func publishCurrentMap(client mqtt.Client) {
	token := client.Publish(currentMapTopic, 0, true, currentMap)
	token.Wait()
	time.Sleep(time.Second)
}
func publishState(client mqtt.Client, status string) {
	token := client.Publish(stateTopic, 0, true, status)
	token.Wait()
	time.Sleep(time.Second)
}

func subCommandTopic(client mqtt.Client) {
	topic := commandTopic
	token := client.Subscribe(topic, 1, messageCommandTopicHandler)
	token.Wait()
	log.Printf("Subscribed to topic: %s", topic)
}

// Subscribes to the current map topic to load the last map loaded
func subCurrentMapTopic(client mqtt.Client) {
	topic := currentMapTopic
	token := client.Subscribe(topic, 1, messageStateTopicHandler)
	token.Wait()
	log.Printf("Subscribed to topic: %s", topic)
}

// Switches the map and reboots the robot
func changeMap(client mqtt.Client, newMap string) {
	publishState(client, "changing_map")
	log.Printf("Changing map from %s to %s\n", currentMap, newMap)

	err := util.RotateFile(roationKeepMaps, fmt.Sprintf("%s/%s", maploaderDir, currentMap), "tar")

	log.Println("saving current map")
	err = tar.Tar(fmt.Sprintf("%s/%s.tar.gz", maploaderDir, currentMap), robot.CurrentRobot.MapFilesAndFolders()...)
	checkAndHandleErrorWithMqtt(err, client)

	log.Println("stopping processes")
	robot.StopProcesses()

	log.Println("removing current map files")
	util.RemoveDirContents(robot.CurrentRobot.MapFolders...)

	mapFileToLoadMatches, err := filepath.Glob(fmt.Sprintf("%s/%s.tar*", maploaderDir, newMap))
	checkAndHandleErrorWithMqtt(err, client)

	if len(mapFileToLoadMatches) == 0 {
		log.Println("requested map does not exist")
	} else {
		log.Println("restoring map from archive")
		err = tar.Untar(mapFileToLoadMatches[0], "/")
		checkAndHandleErrorWithMqtt(err, client)
	}

	log.Println("map change complete, syncing files")
	robot.ExcuteCmd("sync")

	log.Println("restarting processes")
	robot.StartProcesses()

	currentMap = newMap
	publishCurrentMap(client)
	publishState(client, "idle")

}

func checkAndHandleErrorWithMqtt(err error, client mqtt.Client) {
	if err != nil {
		util.CheckAndHandleError(err)
		publishState(client, "error")
	}
}
