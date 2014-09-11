package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kylelemons/goat/termios"
)

var CTRL_C = []byte{3}
var CTRL_D = []byte{4}
var CTRL_Z = []byte{0x1a}

func main() {

	tio, err := termios.NewTermSettings(syscall.Stdin)
	defer tio.Reset()
	if err != nil {
		log.Fatal(err)
	}

	width, height, err := tio.GetSize()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Console size: %dx%d\n", width, height)

	tio.Raw()
	if err != nil {
		log.Fatal(err)
	}

	// Restore tty on kill/interrupt signals

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		os.Interrupt,
		os.Kill,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGQUIT)

	go func() {
		sig := <-c
		// fmt.Printf("%v", sig)
		tio.Reset()
		signal.Stop(c)
		syscall.Kill(os.Getpid(), sig.(syscall.Signal))
	}()

	var b []byte = make([]byte, 5)
	for {
		n, err := os.Stdin.Read(b)
		if err != nil {
			log.Fatal(err)
			break
		}
		if bytes.Equal(b[0:n], CTRL_C) || bytes.Equal(b[0:n], CTRL_D) || bytes.Equal(b[0:n], CTRL_Z) {
			break
		}
		if n > 0 {
			fmt.Printf("% x\n\r", b[0:n])
		}
	}
}
