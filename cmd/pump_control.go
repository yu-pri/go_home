package cmd

import (
	"os"
	"os/signal"
	"syscall"
)

type PumpControllerConfig struct {
	brokerUri              string
	reverseTempSensorTopic string
	pumpRelayTopic         string
}

func main() {
	client := SetupMqttConnection()
	defer client.Disconnect(1000)

	client.Subscribe("/sensors/temp/01", 1, HandleTemperaturePayload)
	defer client.Unsubscribe("/sensors/temp/01")

	<-setupStopSignal()
}

func setupStopSignal() chan bool {
	osSignals := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)

	go func() {
		<-osSignals
		done <- true
	}()

	return done
}
