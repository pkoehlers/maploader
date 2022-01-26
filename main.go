package main

import (
	"errors"
	"fmt"
	"log"
	"maploader/config"
	"maploader/tar"
	"maploader/util"
	"os"
	"os/exec"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const stateTopic string = "valetudo/maploader/map"
const commandTopic string = stateTopic + "/set"

var maploaderDir string
var currentMap string

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	if strings.HasSuffix(msg.Topic(), "set") {
		changeMap(client, string(msg.Payload()))
	} else if string(msg.Payload()) != currentMap {
		log.Printf("Loaded current map from status topic")
		currentMap = string(msg.Payload())
	}

}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}

func main() {
	currentMap = config.Getenv("DEFAULT_MAP_NAME", "main")
	maploaderDir = config.Getenv("MAPLOADER_DIR", "/data/maploader")

	util.InitLogging(maploaderDir + "/log")
	config.InitConfig(config.Getenv("VALETUDO_CONFIG_PATH", "/data/valetudo_config.json"))

	var broker = config.MqttHost()
	var port = config.MqttPort()
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("maploader")
	opts.SetUsername(config.MqttUsername())
	opts.SetPassword(config.MqttPassword())
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
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

	sub(client)
	subRetain(client)

	for {
		time.Sleep(time.Second)
	}
}

func publish(client mqtt.Client) {
	token := client.Publish(stateTopic, 0, true, currentMap)
	token.Wait()
	time.Sleep(time.Second)
}

func sub(client mqtt.Client) {
	topic := commandTopic
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	log.Printf("Subscribed to topic: %s", topic)
}

// Subscribes to the state topic to load the last map loaded
func subRetain(client mqtt.Client) {
	topic := stateTopic
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	log.Printf("Subscribed to topic: %s", topic)
}

// Switches the map and reboots the robot
func changeMap(client mqtt.Client, newMap string) {
	log.Printf("Changing map from %s to %s\n", currentMap, newMap)

	err := util.RotateFile(5, fmt.Sprintf("%s/%s", maploaderDir, currentMap), "tar")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("saving current map")
	tar.Tar(fmt.Sprintf("%s/%s.tar", maploaderDir, currentMap), "/data/ri", "/data/map", "/data/DivideMap", "/data/config/ava/mult_map.json")

	log.Println("removing current map files")
	util.RemoveDirContents("/data/ri/")
	util.RemoveDirContents("/data/map/")
	util.RemoveDirContents("/data/DivideMap/")

	mapFileToLoad := fmt.Sprintf("%s/%s.tar", maploaderDir, newMap)

	if _, err := os.Stat(mapFileToLoad); errors.Is(err, os.ErrNotExist) {
		log.Println("requested map does not exist")
	} else {
		log.Println("restoring map from archive")
		err = tar.Untar(fmt.Sprintf("%s/%s.tar", maploaderDir, newMap), "/")
	}

	log.Println("map change complete, rebooting robot")
	cmd := exec.Command("sync")

	errSync := cmd.Run()

	if errSync != nil {
		log.Fatal(errSync)
	}

	currentMap = newMap
	publish(client)

	time.Sleep(4 * time.Second)

	cmdReboot := exec.Command("reboot")
	errReboot := cmdReboot.Run()

	if errReboot != nil {
		log.Fatal(errReboot)
	}
}
