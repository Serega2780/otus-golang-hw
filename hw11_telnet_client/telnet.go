package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

const (
	TCP   = "tcp"
	COLON = ":"
	bSize = 1024
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TC struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
	abortCh chan interface{}
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	// Place your code here.
	return &TC{address: address, timeout: timeout, in: in, out: out, abortCh: make(chan interface{})}
}

func (tc *TC) Connect() error {
	server := strings.Split(tc.address, COLON)[0]
	_, _ = fmt.Fprintln(os.Stderr, "Trying "+server+"...")
	dialer := net.Dialer{Timeout: tc.timeout, KeepAlive: time.Duration(5 * float64(time.Second))}
	conn, err := dialer.Dial(TCP, tc.address)
	if err == nil {
		_, _ = fmt.Fprintf(os.Stderr, "Connected to %s.\n", server)
		_, _ = fmt.Fprintln(os.Stderr, "Press 'Ctrl+D or Ctrl+C' to quit")
		tc.conn = conn
	}
	return err
}

func (tc *TC) Close() error {
	_, _ = fmt.Fprintln(os.Stderr, "Closing connection... ")
	if err := tc.conn.Close(); err != nil {
		return err
	}
	_, _ = fmt.Fprintln(os.Stderr, "Exited.")
	return nil
}

func (tc *TC) Send() error {
	buf := [bSize]byte{}
	in, err := tc.in.Read(buf[:])
	if errors.Is(err, io.EOF) {
		close(tc.abortCh)
	} else {
		_, err = io.Copy(tc.conn, bytes.NewReader(buf[0:in]))
		if err != nil {
			return err
		}
	}

	return nil
}

func (tc *TC) Receive() error {
	_, err := io.Copy(tc.out, tc.conn)
	if err != nil {
		return err
	}
	return nil
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
