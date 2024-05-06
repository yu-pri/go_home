package pumpcontrol

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func HandleTemperaturePayload(c mqtt.Client, m mqtt.Message) {
	payload := m.Payload()

	fmt.Printf("received payload : %s", payload)
}
