package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var jsonConfigFile string

func Getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func InitConfig(configFile string) {
	jsonConfigFile = configFile
}

func MqttHost() string {
	config := getValetudoConfig()
	return config.Mqtt.Connection.Host
}

func MqttPort() int {
	config := getValetudoConfig()
	return config.Mqtt.Connection.Port
}

func MqttUsername() string {
	config := getValetudoConfig()
	return config.Mqtt.Connection.Authentication.Credentials.Username
}
func MqttPassword() string {
	config := getValetudoConfig()
	return config.Mqtt.Connection.Authentication.Credentials.Password
}

var valetudoConfig ValetudoConfig

func getValetudoConfig() ValetudoConfig {
	if (valetudoConfig == ValetudoConfig{}) {
		configJson, err := ioutil.ReadFile(jsonConfigFile)
		if err != nil {
			log.Print(err)
		}
		json.Unmarshal([]byte(configJson), &valetudoConfig)

	}

	return valetudoConfig

}
