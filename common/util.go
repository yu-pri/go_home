package common

import (
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v2"
)

func ReadParseYamlFile[T interface{}](filePath string) (*T, error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var out T

	err = yaml.Unmarshal(yamlFile, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func SetupStopSignal() chan bool {
	osSignals := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)

	go func() {
		<-osSignals
		done <- true
	}()

	return done
}
