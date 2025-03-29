package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Client struct {
	connection net.Conn
	bufferSize int
}

func NewClient(address string, bufferSize int) (*Client, error) {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return &Client{
		connection: connection,
		bufferSize: bufferSize,
	}, nil
}

func (c *Client) Send(request []byte) ([]byte, error) {
	_, err := c.connection.Write(request)
	if err != nil {
		return nil, err
	}

	response := make([]byte, c.bufferSize)
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
