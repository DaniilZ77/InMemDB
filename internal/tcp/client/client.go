package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

const defaultBufferSize = 1024

type Client struct {
	connection  net.Conn
	readTimeout time.Duration
	bufferSize  int
}

func NewClient(address string, opts ...ClientOption) (*Client, error) {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	client := &Client{connection: connection}
	for _, opt := range opts {
		opt(client)
	}

	if client.bufferSize == 0 {
		client.bufferSize = defaultBufferSize
	}

	return client, nil
}

func (c *Client) Send(request []byte) ([]byte, error) {
	_, err := c.connection.Write(request)
	if err != nil {
		return nil, err
	}

	response := make([]byte, c.bufferSize)
	if c.readTimeout != 0 {
		if err := c.connection.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			return nil, err
		}
	}
	n, err := c.connection.Read(response)
	if err != nil {
		return nil, err
	}

	return response[:n], nil
}

func (c *Client) Close() error {
	return c.connection.Close()
}

func (c *Client) Run() error {
	defer c.Close() // nolint

	stdinReader := bufio.NewReader(os.Stdin)
	request := make([]byte, c.bufferSize)
	for {
		fmt.Print("# ")

		n, err := stdinReader.Read(request)
		if err != nil {
			return err
		}

		response, err := c.Send(request[:n])
		if err != nil {
			return err
		}

		fmt.Println(string(response))
	}
}
