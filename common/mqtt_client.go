package common

import (
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttClientConfig struct {
	passwdFile string `yaml:"passwd_file"`
	brokerUri  string `yaml:"broker_uri"`
	clientId   string `yaml:"client_id"`
	verbose    bool   `yaml:"verbose"`
}

func ParseMqttConfigFile(filePath string) (MqttClientConfig, error) {
	cfg, err := ReadParseYamlFile[MqttClientConfig](filePath)
	return *cfg, err
}

type MqttClientSecrets struct {
	username string
	password string
}

func ReadMqttClientSecretsFromEnv() MqttClientSecrets {
	username := os.Getenv("MQTT_USERNAME")
	password := os.Getenv("MQTT_PASSWORD")

	return MqttClientSecrets{
		username, password,
	}
}

type logger struct{}

func (logger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
func (logger) Println(v ...interface{}) {
	fmt.Println(v...)
}

func SetupMqttConnection(c MqttClientConfig, s MqttClientSecrets) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.SetUsername(s.username)
	opts.SetPassword(s.password)
	opts.AddBroker(c.brokerUri)
	opts.SetClientID(c.clientId)

	logger := logger{}
	mqtt.ERROR = logger
	mqtt.CRITICAL = logger
	mqtt.WARN = logger

	if c.verbose {
		mqtt.DEBUG = logger
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}
