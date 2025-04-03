package tests

import (
	"net"
	"os"
	"testing"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	masterAddr = "localhost:3223"
	slaveAddr  = "localhost:3224"
)

func sendRequest(t *testing.T, connection net.Conn, request string) string {
	response := make([]byte, 1024)
	_, err := common.Write(connection, []byte(request))
	require.NoError(t, err)
	n, err := common.Read(connection, response)
	require.NoError(t, err)
	return string(response[:n])
}

func TestMasterApi(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping api test in short mode")
	}

	t.Cleanup(func() {
		err := os.RemoveAll("./tests/testdata")
		require.NoError(t, err)
	})

	var err error
	connections := make([]net.Conn, 2)
	connections[0], err = net.Dial("tcp", masterAddr)
	require.NoError(t, err)
	defer connections[0].Close() // nolint

	connections[1], err = net.Dial("tcp", slaveAddr)
	require.NoError(t, err)
	defer connections[1].Close() // nolint

	response := sendRequest(t, connections[0], "set name Daniil")
	assert.Equal(t, "OK", response)

	response = sendRequest(t, connections[0], "get name")
	assert.Equal(t, "Daniil", response)

	response = sendRequest(t, connections[0], "del name")
	assert.Equal(t, "OK", response)

	response = sendRequest(t, connections[0], "get name")
	assert.Equal(t, "NIL", response)
}

func TestSlaveApi(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping api test in short mode")
	}

	t.Cleanup(func() {
		err := os.RemoveAll("./tests/testdata")
		require.NoError(t, err)
	})

	var err error
	connections := make([]net.Conn, 2)
	connections[0], err = net.Dial("tcp", masterAddr)
	require.NoError(t, err)
	defer connections[0].Close() // nolint

	connections[1], err = net.Dial("tcp", slaveAddr)
	require.NoError(t, err)
	defer connections[1].Close() // nolint

	const iterationsNumber = 100
	requests := []string{
		"set name Daniil",
		"set age 22",
		"set university MIT",
	}
	for i := range iterationsNumber {
		response := sendRequest(t, connections[0], requests[i%len(requests)])
		assert.Equal(t, "OK", response)
	}

	time.Sleep(100 * time.Millisecond)

	response := sendRequest(t, connections[1], "get name")
	assert.Equal(t, "Daniil", response)

	response = sendRequest(t, connections[1], "get age")
	assert.Equal(t, "22", response)

	response = sendRequest(t, connections[1], "get university")
	assert.Equal(t, "MIT", response)
}
