package main

import (
	"fmt"
	"os"

	"github.com/joshbarrass/TerrariaServerWrapper/internal"
	"github.com/sirupsen/logrus"
)

var mainCommand = "./TerrariaServer.bin.x86_64"

func main() {
	command := append([]string{mainCommand}, os.Args...)
	server, err := internal.NewServer(command)
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
