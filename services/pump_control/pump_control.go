package pumpcontrol

import (
	"fmt"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/yu-pri/go_home/common"
)

type PumpController struct {
	mqttClient               mqtt.Client
	isRelayEnabled           bool
	isRelayManuallyEnabled   bool
	latestReverseTemperature int
	config                   PumpControllerConfig
}

func (c *PumpController) ToggleRelayIfNeeded(isEnabled bool) {
	if isEnabled == c.isRelayEnabled {
		return
	}

	c.mqttClient.Publish(c.config.PumpRelayTopic, common.QosAtLeastOnce, true, isEnabled)
}

func (pumpController *PumpController) ReverseTemperatureHandler() mqtt.MessageHandler {
	return func(c mqtt.Client, m mqtt.Message) {
		payload := m.Payload()
		fmt.Printf("received payload : %s", payload)

		temperatureCelsius, _ := strconv.Atoi(string(payload))
		pumpController.latestReverseTemperature = temperatureCelsius

		relayShouldBeEnabled := temperatureCelsius > pumpController.config.PumpOnTemperature
		pumpController.ToggleRelayIfNeeded(relayShouldBeEnabled)
	}
}

func (pumpController *PumpController) DesiredStateHandler(config PumpControllerConfig) mqtt.MessageHandler {
	return func(c mqtt.Client, m mqtt.Message) {
		payload := m.Payload()

		desiredState := string(payload)

		if desiredState == PumpDesiredStateOn {
			pumpController.ToggleRelayIfNeeded(true)
			// todo set manual mode to prevent automatic relay toggles
			return
		}

		if desiredState == PumpDesiredStateOff {
			pumpController.ToggleRelayIfNeeded(false)
			return
		}

		if desiredState == PumpDesiredStateAuto {
			// TODO: set relay state based on latestReverseTemperature
			return
		}
	}
}
