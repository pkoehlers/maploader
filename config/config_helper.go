package config

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
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
func RotationKeepMaps() int {
	rotationKeepMaps, err := strconv.Atoi(Getenv("ROTATION_KEEP_MAPS", "5"))
	if err != nil {
		rotationKeepMaps = 5
	}
	return rotationKeepMaps
}

var valetudoConfig ValetudoConfig

func getValetudoConfig() ValetudoConfig {
	if (valetudoConfig == ValetudoConfig{}) {
		configJson, err := os.ReadFile(jsonConfigFile)
		if err != nil {
			log.Print(err)
		}
		json.Unmarshal([]byte(configJson), &valetudoConfig)

	}

	return valetudoConfig

}
