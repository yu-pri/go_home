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
	inManualMode             bool
	latestReverseTemperature *int
	config                   PumpControllerConfig
}

func NewPumpController(mqttClient mqtt.Client, config PumpControllerConfig) PumpController {
	return PumpController{
		mqttClient:               mqttClient,
		config:                   config,
		inManualMode:             true, // wait until broker tells we're in auto mode
		isRelayEnabled:           false,
		latestReverseTemperature: nil, // TODO save a history log; restore if latest reading isn't far off in the past
	}
}

func (c *PumpController) SetupSubscriptions() {
	c.mqttClient.Subscribe(c.config.PumpDesiredStateTopic, common.QosAtLeastOnce, c.DesiredStateHandler())
	c.mqttClient.Subscribe(c.config.ReverseTempSensorTopic, common.QosAtLeastOnce, c.ReverseTemperatureHandler())
}

func (pumpController *PumpController) DesiredStateHandler() mqtt.MessageHandler {
	return func(c mqtt.Client, m mqtt.Message) {
		payload := m.Payload()
		fmt.Printf("received payload : %s", payload)

		desiredState := string(payload)

		if PumpDesiredStateOn == desiredState {
			pumpController.inManualMode = true
			pumpController.toggleRelayIfNeeded(true)
			return
		}

		if PumpDesiredStateOff == desiredState {
			pumpController.inManualMode = true
			pumpController.toggleRelayIfNeeded(false)
			return
		}

		if PumpDesiredStateAuto == desiredState {
			pumpController.inManualMode = false
			pumpController.autoSetRelayState(pumpController.latestReverseTemperature)
			return
		}

		fmt.Printf("Bad payload for desired pump state topic: %s", payload)
	}
}

func (pumpController *PumpController) ReverseTemperatureHandler() mqtt.MessageHandler {
	return func(c mqtt.Client, m mqtt.Message) {
		payload := m.Payload()
		fmt.Printf("received payload : %s", payload)

		temperatureCelsius, err := strconv.Atoi(string(payload))
		if err != nil {
			fmt.Printf("Bad payload for reverse temperature topic: %s, error: %s", payload, err)
			return
		}

		pumpController.latestReverseTemperature = &temperatureCelsius

		if pumpController.inManualMode {
			return
		}

		pumpController.autoSetRelayState(&temperatureCelsius)
	}
}

func (pumpController *PumpController) autoSetRelayState(currentReverseTemperature *int) {
	if currentReverseTemperature == nil {
		return
	}
	relayShouldBeEnabled := *currentReverseTemperature > pumpController.config.PumpOnTemperature
	pumpController.toggleRelayIfNeeded(relayShouldBeEnabled)
}

func (c *PumpController) toggleRelayIfNeeded(isEnabled bool) {
	if isEnabled == c.isRelayEnabled {
		return
	}

	token := c.mqttClient.Publish(c.config.PumpRelayTopic, common.QosAtLeastOnce, true, isEnabled)

	if token.Wait() && token.Error() != nil {
		fmt.Printf("Failed to communicate with relay service: %s", token.Error().Error())
		return
	}

	c.isRelayEnabled = isEnabled
}
