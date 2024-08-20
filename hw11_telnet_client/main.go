package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	TimeoutFlag  = "timeout"
	TimeoutValue = 10
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Duration(TimeoutValue*float64(time.Second)), "timeout to connect to server")
}

func main() {
	// Place your code here,
	// P.S. Do not rush to throw context down, think if it is useful with blocking operation?
	stop := make(chan interface{})
	var host, port string
	arguments := os.Args
	flag.Parse()

	if (isFlagPassed(TimeoutFlag) && len(arguments) < 4) || (!isFlagPassed(TimeoutFlag) && len(arguments) < 3) {
		fmt.Println("Please provide host port as separate arguments")
		return
	}
	if isFlagPassed(TimeoutFlag) {
		host = os.Args[2]
		port = os.Args[3]
	} else {
		host = os.Args[1]
		port = os.Args[2]
	}

	tc := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	if err := tc.Connect(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	go func(tc TelnetClient) {
		for {
			_ = tc.Receive()
		}
	}(tc)

	go func(tc TelnetClient) {
		for {
			err := tc.Send()
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
		}
	}(tc)

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGURG)
		sig := <-ch
		if sig == syscall.SIGURG {
			_, _ = fmt.Fprintln(os.Stderr, "...EOF")
			if err := tc.Close(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
		}
		os.Exit(1)
	}()

	<-stop
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
