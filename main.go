package main

import (
	"fmt"
	"log"
	"maploader/config"
	"maploader/tar"
	"maploader/util"
	"os"
	"os/exec"
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
	err = tar.Tar(fmt.Sprintf("%s/%s.tar.gz", maploaderDir, currentMap), "/data/ri", "/data/map", "/data/DivideMap", "/data/config/ava/mult_map.json")
	checkAndHandleErrorWithMqtt(err, client)

	log.Println("stopping processes")
	stopProcesses()

	log.Println("removing current map files")
	util.RemoveDirContents("/data/ri/")
	util.RemoveDirContents("/data/map/")
	util.RemoveDirContents("/data/DivideMap/")

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
	excuteCmd("sync")

	log.Println("restarting processes")
	startProcesses()

	currentMap = newMap
	publishCurrentMap(client)
	publishState(client, "idle")

}

func stopProcesses() {
	excuteCmd("killall", "-9", "valetudo")
	excuteCmd("sh", "/etc/rc.d/miio.sh", "stop")
	excuteCmd("killall", "-9", "ava")
}

func startProcesses() {
	excuteCmd("sh", "/etc/rc.d/miio.sh")
	excuteCmd("sh", "/etc/rc.d/ava.sh")
	startValetudo()
}

func startValetudo() {
	devnull, dnerr := os.OpenFile(os.DevNull, os.O_WRONLY, 0755)
	if dnerr != nil {
		util.CheckAndHandleError(dnerr)
	}

	cmd := exec.Command("/data/valetudo")
	cmd.Stdout = devnull
	cmd.Env = os.Environ()
	valetudoConfigEnv := "VALETUDO_CONFIG_PATH=" + config.Getenv("VALETUDO_CONFIG_PATH", "/data/valetudo_config.json")
	cmd.Env = append(cmd.Env, valetudoConfigEnv)
	err := cmd.Start()

	if err != nil {
		util.CheckAndHandleError(err)
	} else {
		cmd.Process.Release()
	}
}

func excuteCmd(cmdStr string, cmdArgs ...string) {

	cmd := exec.Command(cmdStr, cmdArgs...)
	err := cmd.Run()

	if err != nil {
		util.CheckAndHandleError(err)
	}
}

func checkAndHandleErrorWithMqtt(err error, client mqtt.Client) {
	if err != nil {
		util.CheckAndHandleError(err)
		publishState(client, "error")
	}
}
