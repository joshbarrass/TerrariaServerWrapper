package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joshbarrass/TerrariaServerWrapper/internal"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Configuration struct {
	AutosaveTime string `envconfig:"AUTOSAVE_TIME" default:"5m"`
	Executable   string `envconfig:"SERVER_EXECUTABLE" default:"./TerrariaServer.bin.x86_64"`
}

func main() {
	var config Configuration
	err := envconfig.Process("", &config)
	if err != nil {
		logrus.Fatalf("Failed to parse envconfig: %s", err)
	}

	autosaveTime, err := time.ParseDuration(config.AutosaveTime)
	if err != nil {
		logrus.Fatalf("Failed to parse autosave time: %s", err)
	}

	command := append([]string{config.Executable}, os.Args[1:]...)
	server, err := internal.NewServer(command, autosaveTime)
	if err != nil {
		logrus.Fatalf("An error occured starting the server: %s", err)
	}
	err = server.Start()
	fmt.Println()
	if err != nil {
		logrus.Errorf("An error occured in the server: %s", err)
	} else {
		logrus.Infof("exited with status code %d", server.GetExitCode())
	}
}
