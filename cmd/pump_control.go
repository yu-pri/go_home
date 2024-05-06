package cmd

import (
	"github.com/yu-pri/go_home/common"
	"github.com/yu-pri/go_home/services/pump_control"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	client := common.SetupMqttConnection()
	defer client.Disconnect(1000)

	client.Subscribe("/sensors/temp/01", 1, pumpcontrol.HandleTemperaturePayload)
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
