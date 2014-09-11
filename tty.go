package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func stty(args ...string) string {
	// fmt.Printf("%v\n\r", args)
	args = append([]string{"-f", "/dev/tty"}, args...)
	out, err := exec.Command("stty", args...).Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

var CTRL_C = []byte{3}
var CTRL_D = []byte{4}
var CTRL_Z = []byte{0x1a}

func main() {
	size := stty("size")
	fmt.Printf("Console size: %s\n", strings.Trim(size, "\n"))

	// Restore stty on end

	restore := stty("-g")
	defer stty(restore)

	// Restore tty on kill/interupt signals

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
		<-c
		// sig := <-c
		// fmt.Printf("%v", sig)
		stty(restore)
		os.Exit(1)
	}()

	// Put the tty in raw mode

	stty("raw")

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
