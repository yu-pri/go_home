package main

import (
	"fmt"
	"os"

	"github.com/yu-pri/go_home/common"
	"github.com/yu-pri/go_home/services/pump_control"
)

func main() {
	mqttConfigPath := os.Args[1]
	pumpConfigPath := os.Args[2]
	fmt.Printf("mqtt config: %s, pump config: %s", mqttConfigPath, pumpConfigPath)

	mqttSecrets := common.ReadMqttClientSecretsFromEnv()
	mqttConfig, err := common.ParseMqttConfigFile(mqttConfigPath)

	if err != nil {
		fmt.Printf("Error parsing configs: %s", err)
		return
	}
	pumpConfig, err := pumpcontrol.ParseConfig(pumpConfigPath)
	if err != nil {
		fmt.Printf("Error parsing configs: %s", err)
		return
	}
	client := common.SetupMqttConnection(*mqttConfig, mqttSecrets)
	defer client.Disconnect(1000)

	pumpController := pumpcontrol.NewPumpController(client, pumpConfig)
	pumpController.SetupSubscriptions()

	<-common.SetupStopSignal()
}
