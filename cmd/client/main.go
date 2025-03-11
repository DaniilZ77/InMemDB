package main

import (
	"flag"

	"github.com/DaniilZ77/InMemDB/internal/tcp/client"
)

func main() {
	var address string
	flag.StringVar(&address, "address", "127.0.0.1:3223", "server address")

	flag.Parse()

	client := client.NewClient(address)
	if err := client.Run(); err != nil {
		panic(err)
	}
}
