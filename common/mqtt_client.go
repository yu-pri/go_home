package common

import (
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttClientConfig struct {
	BrokerUri string `yaml:"broker_uri"`
	Verbose   bool   `yaml:"verbose"`
}

func ParseMqttConfigFile(filePath string) (*MqttClientConfig, error) {
	out, err := ReadParseYamlFile[MqttClientConfig](filePath)
	fmt.Printf("\nmqtt cfg: %s, %s, %s\n", out.BrokerUri, out.ClientId, out.Verbose)
	return out, err
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

func SetupMqttConnection(c MqttClientConfig, s MqttClientSecrets, clientId string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	fmt.Printf("broker uri : %s", c.BrokerUri)
	opts.SetUsername(s.username)
	opts.SetPassword(s.password)
	opts.AddBroker(c.BrokerUri)
	opts.SetClientID(clientId)

	logger := logger{}
	mqtt.ERROR = logger
	mqtt.CRITICAL = logger
	mqtt.WARN = logger

	if c.Verbose {
		mqtt.DEBUG = logger
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}
