package main

import (
	"bufio"
	"context"
	"io"
	"net"
	"time"

	"github.com/pkg/errors"
)

type TelnetClient interface {
	Connect(ctx context.Context) error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (client *Client) Close() error {
	return client.conn.Close()
}

func (client *Client) Connect(ctx context.Context) error {
	var err error
	client.conn, err = net.DialTimeout("tcp", client.address, client.timeout)
	if err != nil {
		return errors.Wrap(err, "error connect")
	}

	go func() {
		<-ctx.Done()
		client.conn.Close()
	}()
	return nil
}

func (client *Client) Send() error {
	r := bufio.NewReader(client.in)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return errors.Wrap(err, "error read")
		}

		b = append(b, '\n')
		if _, err = client.conn.Write(b); err != nil {
			return errors.Wrap(err, "error write")
		}
	}
}

func (client *Client) Receive() error {
	r := bufio.NewReader(client.conn)
	for {
		b, _, err := r.ReadLine()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return errors.Wrap(err, "error read")
		}

		b = append(b, '\n')
		_, err = client.out.Write(b)
		if err != nil {
			return errors.Wrap(err, "error write")
		}
	}
}
