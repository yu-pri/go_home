package pumpcontrol

import "github.com/yu-pri/go_home/common"

type PumpControllerConfig struct {
	ReverseTempSensorTopic string `yaml:"reverse_temp_sensor_topic"`
	PumpRelayTopic         string `yaml:"pump_relay_topic"`
	PumpDesiredStateTopic  string `yaml:"pump_desired_state_topic"`
	PumpOnTemperature      int    `yaml:"pump_on_temperature"`
}

func ParseConfig(configPath string) (PumpControllerConfig, error) {
	config, err := common.ReadParseYamlFile[PumpControllerConfig](configPath)
	return *config, err
}
