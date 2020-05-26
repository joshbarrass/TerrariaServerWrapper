package internal

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

const exitCommand = "exit\r"

// Server type holds all the os/exec objects for the terraria server
type Server struct {
	Command   *exec.Cmd
	Stdin     io.WriteCloser
	quit      chan struct{}
	lastWrite time.Time
}

// NewServer launches a new Terraria server with a given command
func NewServer(command []string) (*Server, error) {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// cmd.Stdin = os.Stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	quit := make(chan struct{})

	return &Server{
		Command:   cmd,
		Stdin:     stdin,
		quit:      quit,
		lastWrite: time.Now(),
	}, nil
}

func (server *Server) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start the server
	err := server.Command.Start()
	if err != nil {
		return err
	}

	fmt.Println("Server starting...")
	server.ShutdownOnExit()
	server.startInputLoop(ctx)
	server.startAutosaveLoop(ctx)
	server.startSigtermHandler(ctx)

	// wait for exit
	<-server.quit

	return nil
}

func (server *Server) Shutdown() error {
	// tell the server to save and exit
	server.Stdin.Write([]byte(autosaveCommand))
	server.Stdin.Write([]byte(exitCommand))

	server.quit <- struct{}{}
	return nil
}

func (server *Server) ShutdownOnExit() {
	go func() {
		server.Command.Wait()
		server.Shutdown()
	}()
}

func (server *Server) GetExitCode() int {
	return server.Command.ProcessState.ExitCode()
}
