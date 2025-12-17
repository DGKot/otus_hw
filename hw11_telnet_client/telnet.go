package main

import (
	"errors"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

var ErrNotConnect = errors.New("client didn't connect")

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	// Place your code here.
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *telnetClient) Connect() error {
	var err error
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)
	return err
}

func (t *telnetClient) Send() error {
	if t.conn == nil {
		return ErrNotConnect
	}
	_, err := io.Copy(t.conn, t.in)
	if err == nil {
		tcp, ok := t.conn.(*net.TCPConn)
		if ok {
			_ = tcp.CloseWrite()
		}
	}
	return err
}

func (t *telnetClient) Receive() error {
	if t.conn == nil {
		return ErrNotConnect
	}
	_, err := io.Copy(t.out, t.conn)
	return err
}

func (t *telnetClient) Close() error {
	if t.conn == nil {
		return nil
	}
	return t.conn.Close()
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
