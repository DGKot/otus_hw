package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", time.Second*10, "timeout connection")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("Miss port for connect")
	}
	host := args[0]
	port := args[1]

	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		log.Fatal("error connect to server:", err)
	}

	closeCh := make(chan struct{})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		_ = client.Send()
	}()

	go func() {
		client.Receive()
		close(closeCh)
	}()

	select {
	case <-sigCh:
		log.Println("SIGINT. Connection is closed")
	case <-closeCh:
	}

	client.Close()
}
