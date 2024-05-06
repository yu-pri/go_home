package common

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type verboseLogger struct{}

func (verboseLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
func (verboseLogger) Println(v ...interface{}) {
	fmt.Println(v...)
}

func SetupMqttConnection() mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.SetUsername("yupri")
	opts.SetPassword("test")
	opts.AddBroker("localhost:1883")
	opts.SetClientID("pump_controller")
	logger := verboseLogger{}
	mqtt.ERROR = logger
	mqtt.DEBUG = logger
	mqtt.WARN = logger
	mqtt.CRITICAL = logger

	client := mqtt.NewClient(opts)
	status := client.IsConnected()
	fmt.Printf("initial status: %t", status)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}
