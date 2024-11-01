package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

func Getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func MqttHost() string {
	config := getAppConfig()
	return config.Mqtt.Connection.Host
}

func MqttPort() int {
	config := getAppConfig()
	return config.Mqtt.Connection.Port
}

func MqttUsername() string {
	config := getAppConfig()
	return config.Mqtt.Connection.Authentication.Credentials.Username
}
func MqttPassword() string {
	config := getAppConfig()
	return config.Mqtt.Connection.Authentication.Credentials.Password
}
func MqttTLSEnabled() bool {
	config := getAppConfig()
	return config.Mqtt.Connection.TLS.Enabled
}
func MqttTLSCA() string {
	config := getAppConfig()
	return config.Mqtt.Connection.TLS.Ca
}
func MqttIdentifier() string {
	return "wakeonlan"
}

func MqttTopicPrefix() string {

	return "wakeonlan"
}

func LookupMacAddressByName(name string) (string, error) {
	for _, device := range AppConfig.Devices {
		if device.Name == name {
			return device.MacAddress, nil
		}
	}
	return "", errors.New("device not found")
}

var (
	AppConfig Config
	once      sync.Once
)

func getAppConfig() Config {

	once.Do(func() {
		jsonConfigFile := Getenv("CONFIG_PATH", "/config.json")
		configJson, err := os.ReadFile(jsonConfigFile)
		if err != nil {
			log.Print(err)
		}
		json.Unmarshal([]byte(configJson), &AppConfig)

	})
	return AppConfig

}
