package cmd

import (
	"github.com/yu-pri/go_home/common"
	"github.com/yu-pri/go_home/services/pump_control"
)

func main() {
	mqttSecrets := common.ReadMqttClientSecretsFromEnv()
	mqttConfig, _ := common.ParseMqttConfigFile("")
	pumpConfig, _ := pumpcontrol.ParseConfig("")

	client := common.SetupMqttConnection(mqttConfig, mqttSecrets)
	defer client.Disconnect(1000)

	client.Subscribe("/sensors/temp/01", 1, pumpcontrol.HandleTemperaturePayload(pumpConfig))
	defer client.Unsubscribe("/sensors/temp/01")

	<-common.SetupStopSignal()
}
