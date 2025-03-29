package main

import (
	"flag"

	"github.com/DaniilZ77/InMemDB/internal/tcp/client"
)

func main() {
	var address string
	flag.StringVar(&address, "address", "127.0.0.1:3223", "server address")
	var bufferSize int
	flag.IntVar(&bufferSize, "buffer_size", 1024, "buffer size")

	flag.Parse()

	client, err := client.NewClient(address, bufferSize)
	if err != nil {
		panic("failed to init client: " + err.Error())
	}

	if err := client.Run(); err != nil {
		panic(err)
	}
}
