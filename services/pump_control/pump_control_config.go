package pumpcontrol

type PumpControllerConfig struct {
	reverseTempSensorTopic string
	pumpRelayTopic         string
}

func ParseFromFile(configPath string) {
}
