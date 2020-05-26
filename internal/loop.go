package internal

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const autosaveCommand = "save\r"

// startInputLoop begins a goroutine that continuously forwards
// os.stdin to the server's stdin pipe
func (server *Server) startInputLoop(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Exiting read loop")
				return
			default:
				buf := make([]byte, 1024)
				n, _ := os.Stdin.Read(buf)
				if n == 0 {
					continue
				}
				buf = bytes.Trim(buf, "\x00")
				server.Stdin.Write(buf)
				server.lastWrite = time.Now()
			}
		}
	}()
}

// startAutosaveLoop starts a goroutine that waits until there has
// been no activity for a set period of time before sending an
// autosave command
func (server *Server) startAutosaveLoop(ctx context.Context) {
	autosave := make(chan struct{})
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Exiting autosave loop")
				return
			case <-autosave:
				if time.Since(server.lastWrite) < server.AutosaveTime {
					continue
				}
				fmt.Println("Autosaving...")
				server.Stdin.Write([]byte(autosaveCommand))
			}
		}
	}()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(server.AutosaveTime)
				autosave <- struct{}{}
			}
		}
	}()
}

// sigtermHandler starts a goroutine which waits for a SIGTERM and
// safely shuts down the server when it receives one
func (server *Server) startSigtermHandler(ctx context.Context) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM)
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-sigChan:
			fmt.Println("Shutting down safely...")
			server.Shutdown()
			return
		}
	}()
}
