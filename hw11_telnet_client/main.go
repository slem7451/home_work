package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "Timeout connect")
	flag.Parse()

	if flag.NArg() < 2 {
		log.Fatal("Not enough arguments in command")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)

	address := host + ":" + port

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Println(err)
		}
	}()

	fmt.Fprintln(os.Stderr, "Connection is openned")

	go func() {
		if err := client.Send(); err != nil {
			log.Println(err)
		}
	}()

	ctxServer, cancelServer := context.WithCancel(context.Background())
	defer cancelServer()

	go func() {
		if err := client.Receive(); err != nil {
			var netErr net.Error

			if errors.As(err, &netErr) && netErr.Timeout() {
				fmt.Fprintln(os.Stderr, "Connection is timed out")
			} else {
				fmt.Fprintln(os.Stderr, "Connection is closed")
			}
		}
		cancelServer()
	}()

	ctxClient, cancelClient := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelClient()

	select {
	case <-ctxServer.Done():
		fmt.Fprintln(os.Stderr, "Connection is closed by peer")
	case <-ctxClient.Done():
		fmt.Fprintln(os.Stderr, "EOF")
	}
}
