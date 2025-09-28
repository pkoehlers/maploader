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
func MqttTLSEnabled() bool {
	config := getValetudoConfig()
	return config.Mqtt.Connection.TLS.Enabled
}
func MqttTLSCA() string {
	config := getValetudoConfig()
	return config.Mqtt.Connection.TLS.Ca
}
func MqttTLSIgnoreCertificateErrors() bool {
	config := getValetudoConfig()
	return config.Mqtt.Connection.TLS.IgnoreCertificateErrors
}
func MqttIdentifier() string {
	config := getValetudoConfig()
	if config.Mqtt.Identity.Identifier != "" {
		return config.Mqtt.Identity.Identifier
	}
	hostname, err := os.Hostname()
	if err != nil {
		return "unknownrobot"
	}
	return hostname
}

func MqttTopicPrefix() string {
	config := getValetudoConfig()
	if config.Mqtt.Customizations.TopicPrefix != "" {
		return config.Mqtt.Customizations.TopicPrefix
	}

	return "valetudo"
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
