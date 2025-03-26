package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	conn    net.Conn
	bufSize int
}

func NewClient(address string, bufSize int) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn, bufSize: bufSize}, nil
}

func (c *Client) Send(req string) (string, error) {
	_, err := fmt.Fprintln(c.conn, req)
	if err != nil {
		return "", err
	}

	resp := make([]byte, c.bufSize)
	n, err := c.conn.Read(resp)
	if err != nil {
		return "", err
	}

	return string(resp[:n]), nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Run() error {
	defer c.Close() // nolint

	clientReader := bufio.NewReader(os.Stdin)

	var req, resp string
	var err error
	for {
		fmt.Print("# ")

		if req, err = clientReader.ReadString('\n'); err != nil {
			return err
		}
		req = strings.Trim(req, " \n")
		if req == "exit" {
			break
		}

		if resp, err = c.Send(req); err != nil {
			return err
		}
		fmt.Print(resp)
	}

	return nil
}
