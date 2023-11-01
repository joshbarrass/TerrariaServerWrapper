package internal

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

const exitCommand = "exit\n\r"

// Server type holds all the os/exec objects for the terraria server
type Server struct {
	Command      *exec.Cmd
	Stdin        io.WriteCloser
	AutosaveTime time.Duration
	quit         chan struct{}
	lastWrite    time.Time
	interactive  bool
}

// NewServer launches a new Terraria server with a given command
func NewServer(command []string, autosaveTime time.Duration, interactive bool) (*Server, error) {
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
		Command:      cmd,
		Stdin:        stdin,
		AutosaveTime: autosaveTime,
		quit:         quit,
		lastWrite:    time.Now(),
		interactive:  interactive,
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
	// server.ShutdownOnExit()
	if server.interactive {
		server.startInputLoop(ctx)
	}
	server.startAutosaveLoop(ctx)
	server.startSigtermHandler(ctx)

	// wait for exit
	<-server.quit
	// cancel()
	server.Command.Wait()

	return nil
}

func (server *Server) Shutdown() error {
	// tell the server to save and exit
	// server.Stdin.Write([]byte(autosaveCommand))
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
