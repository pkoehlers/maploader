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

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	if strings.HasSuffix(msg.Topic(), "set") {
		changeMap(client, string(msg.Payload()))
	} else if string(msg.Payload()) != currentMap {
		log.Printf("Loaded current map from status topic")
		currentMap = string(msg.Payload())
	}

}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

var currentMap string = "main"

func main() {

	util.InitLogging()
	os.Setenv("VALETUDO_CONFIG_PATH", "valetudo_config.json")
	config.InitConfig(os.Getenv("VALETUDO_CONFIG_PATH"))

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
	fmt.Printf("Subscribed to topic: %s", topic)
}

// Subscribes to the state topic to load the last map loaded
func subRetain(client mqtt.Client) {
	topic := stateTopic
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", topic)
}

// Switches the map and reboots the robot
func changeMap(client mqtt.Client, newMap string) {
	log.Printf("Changing map from %s to %s\n", currentMap, newMap)

	err := util.RotateFile(5, fmt.Sprintf("/data/maploader/%s", currentMap), "tar")
	if err != nil {
		log.Fatal(err)
	}

	tar.Tar(fmt.Sprintf("/data/maploader/%s.tar", currentMap), "/data/ri", "/data/map", "/data/DivideMap", "/data/config/ava/mult_map.json")

	os.RemoveAll("/data/ri/")
	os.MkdirAll("/data/ri/", 0755)

	os.RemoveAll("/data/map/")
	os.MkdirAll("/data/map/", 0755)

	os.RemoveAll("/data/DivideMap/")
	os.MkdirAll("/data/DivideMap/", 0755)
	mapFileToLoad := fmt.Sprintf("/data/maploader/%s.tar", newMap)

	if _, err := os.Stat(mapFileToLoad); errors.Is(err, os.ErrNotExist) {
		log.Println("requested map does not exist")
	} else {
		err = tar.Untar(fmt.Sprintf("/data/maploader/%s.tar", newMap), "/")
	}

	fmt.Println("map change complete, rebooting robot")
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
