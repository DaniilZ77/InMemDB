package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	newLine       = '\n'
	requestCutset = " \n"
)

type Client struct {
	address string
}

func NewClient(address string) *Client {
	return &Client{address: address}
}

func (c *Client) Run() error {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return err
	}
	defer conn.Close()

	user := bufio.NewReader(os.Stdin)
	server := bufio.NewReader(conn)

	var req, resp string
	for {
		fmt.Print("# ")

		if req, err = user.ReadString(newLine); err != nil {
			return err
		}

		req = strings.Trim(req, requestCutset)
		if req == "exit" {
			break
		}

		if _, err = fmt.Fprintln(conn, req); err != nil {
			return err
		}

		if resp, err = server.ReadString(newLine); err != nil {
			return err
		}

		fmt.Print(resp)
	}

	return nil
}
