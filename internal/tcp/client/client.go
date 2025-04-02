package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/common"
)

const defaultBufferSize = 1024

type Client struct {
	connection  net.Conn
	idleTimeout time.Duration
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
	if c.idleTimeout != 0 {
		if err := c.connection.SetWriteDeadline(time.Now().Add(c.idleTimeout)); err != nil {
			return nil, err
		}
	}
	_, err := common.Write(c.connection, request)
	if err != nil {
		return nil, err
	}

	response := make([]byte, c.bufferSize)
	if c.idleTimeout != 0 {
		if err := c.connection.SetReadDeadline(time.Now().Add(c.idleTimeout)); err != nil {
			return nil, err
		}
	}
	n, err := common.Read(c.connection, response)
	if err != nil && err != io.EOF {
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
	for {
		fmt.Print("# ")

		request, err := stdinReader.ReadBytes('\n')
		if err != nil {
			return err
		}

		response, err := c.Send(request)
		if err != nil {
			return err
		}

		fmt.Println(string(response))
	}
}
