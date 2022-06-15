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

const stateTopic string = "valetudo/maploader/map"
const commandTopic string = stateTopic + "/set"

var maploaderDir string
var currentMap string

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
	subStateTopic(client)
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

	for {
		time.Sleep(time.Second)
	}
}

func publishState(client mqtt.Client) {
	token := client.Publish(stateTopic, 0, true, currentMap)
	token.Wait()
	time.Sleep(time.Second)
}

func subCommandTopic(client mqtt.Client) {
	topic := commandTopic
	token := client.Subscribe(topic, 1, messageCommandTopicHandler)
	token.Wait()
	log.Printf("Subscribed to topic: %s", topic)
}

// Subscribes to the state topic to load the last map loaded
func subStateTopic(client mqtt.Client) {
	topic := stateTopic
	token := client.Subscribe(topic, 1, messageStateTopicHandler)
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
	tar.Tar(fmt.Sprintf("%s/%s.tar.gz", maploaderDir, currentMap), "/data/ri", "/data/map", "/data/DivideMap", "/data/config/ava/mult_map.json")

	log.Println("stopping processes")
	stopProcesses()

	log.Println("removing current map files")
	util.RemoveDirContents("/data/ri/")
	util.RemoveDirContents("/data/map/")
	util.RemoveDirContents("/data/DivideMap/")

	mapFileToLoadMatches, err := filepath.Glob(fmt.Sprintf("%s/%s.tar*", maploaderDir, newMap))
	if err != nil {
		log.Fatal(err)
	}
	if len(mapFileToLoadMatches) == 0 {
		log.Println("requested map does not exist")
	} else {
		log.Println("restoring map from archive")
		err = tar.Untar(mapFileToLoadMatches[0], "/")
	}

	log.Println("map change complete, syncing files")
	cmd := exec.Command("sync")

	errSync := cmd.Run()

	if errSync != nil {
		log.Fatal(errSync)
	}

	log.Println("restarting processes")
	startProcesses()

	currentMap = newMap
	publishState(client)

}

func stopProcesses() {
	excuteCmdLogOnError("killall", "-9", "valetudo")
	excuteCmdLogOnError("sh", "/etc/rc.d/miio.sh", "stop")
	excuteCmdLogOnError("killall", "-9", "ava")
}

func startProcesses() {
	excuteCmdLogOnError("sh", "/etc/rc.d/miio.sh")
	excuteCmdLogOnError("sh", "/etc/rc.d/ava.sh")
	startValetudo()
}

func startValetudo() {
	devnull, dnerr := os.OpenFile(os.DevNull, os.O_WRONLY, 0755)
	if dnerr != nil {
		panic(dnerr)
	}

	cmd := exec.Command("/data/valetudo")
	cmd.Stdout = devnull
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "VALETUDO_CONFIG_PATH=/data/valetudo_config.json")
	err := cmd.Start()

	if err != nil {
		log.Fatal(err)
	} else {
		cmd.Process.Release()
	}
}

func excuteCmdLogOnError(cmdStr string, cmdArgs ...string) {

	cmd := exec.Command(cmdStr, cmdArgs...)
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}
