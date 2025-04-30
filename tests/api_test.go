package tests

import (
	"net"
	"os"
	"sync"
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

	connection, err := net.Dial("tcp", masterAddr)
	require.NoError(t, err)
	defer connection.Close() // nolint

	assert.Equal(t, "OK", sendRequest(t, connection, "set name Daniil"))
	assert.Equal(t, "Daniil", sendRequest(t, connection, "get name"))
	assert.Equal(t, "OK", sendRequest(t, connection, "del name"))
	assert.Equal(t, "NIL", sendRequest(t, connection, "get name"))
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
		assert.Equal(t, "OK", sendRequest(t, connections[0], requests[i%len(requests)]))
	}

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, "Daniil", sendRequest(t, connections[1], "get name"))
	assert.Equal(t, "22", sendRequest(t, connections[1], "get age"))
	assert.Equal(t, "MIT", sendRequest(t, connections[1], "get university"))
}

func TestTransactionBasic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping api test in short mode")
	}

	t.Cleanup(func() {
		err := os.RemoveAll("./tests/testdata")
		require.NoError(t, err)
	})

	connection, err := net.Dial("tcp", masterAddr)
	require.NoError(t, err)
	defer connection.Close() // nolint

	assert.Equal(t, "OK", sendRequest(t, connection, "begin"))
	assert.Equal(t, "OK", sendRequest(t, connection, "set name Daniil"))
	assert.Equal(t, "Daniil", sendRequest(t, connection, "get name"))
	assert.Equal(t, "OK", sendRequest(t, connection, "del name"))
	assert.Equal(t, "NIL", sendRequest(t, connection, "get name"))
	assert.Equal(t, "OK", sendRequest(t, connection, "set age 22"))
	assert.Equal(t, "OK", sendRequest(t, connection, "commit"))
	assert.Equal(t, "22", sendRequest(t, connection, "get age"))
	assert.Equal(t, "OK", sendRequest(t, connection, "begin"))
	assert.Equal(t, "OK", sendRequest(t, connection, "del age"))
	assert.Equal(t, "OK", sendRequest(t, connection, "rollback"))
	assert.Equal(t, "22", sendRequest(t, connection, "get age"))
}

func TestTransactionDirtyRead(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping api test in short mode")
	}

	t.Cleanup(func() {
		err := os.RemoveAll("./tests/testdata")
		require.NoError(t, err)
	})

	connection, err := net.Dial("tcp", masterAddr)
	require.NoError(t, err)
	defer connection.Close() // nolint

	assert.Equal(t, "OK", sendRequest(t, connection, "begin"))
	assert.Equal(t, "OK", sendRequest(t, connection, "set a b"))
	assert.Equal(t, "OK", sendRequest(t, connection, "set c d"))

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		connection, err := net.Dial("tcp", masterAddr)
		require.NoError(t, err)
		defer connection.Close() // nolint
		assert.Equal(t, "OK", sendRequest(t, connection, "begin"))
		assert.Equal(t, "NIL", sendRequest(t, connection, "get a"))
		assert.Equal(t, "NIL", sendRequest(t, connection, "get c"))
		assert.Equal(t, "OK", sendRequest(t, connection, "commit"))
	}()
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, "OK", sendRequest(t, connection, "rollback"))
	wg.Wait()
}

func TestTransactionNonRepeatableRead(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping api test in short mode")
	}

	t.Cleanup(func() {
		err := os.RemoveAll("./tests/testdata")
		require.NoError(t, err)
	})

	connection, err := net.Dial("tcp", masterAddr)
	require.NoError(t, err)
	defer connection.Close() // nolint

	assert.Equal(t, "OK", sendRequest(t, connection, "begin"))
	assert.Equal(t, "OK", sendRequest(t, connection, "set e f"))
	go func() {
		connection, err := net.Dial("tcp", masterAddr)
		require.NoError(t, err)
		defer connection.Close() // nolint

		assert.Equal(t, "OK", sendRequest(t, connection, "begin"))
		assert.Equal(t, "OK", sendRequest(t, connection, "set e g"))
		assert.Equal(t, "OK", sendRequest(t, connection, "commit"))
	}()
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, "f", sendRequest(t, connection, "get e"))
	assert.Contains(t, sendRequest(t, connection, "commit"), "ERROR")
}

func TestTransactionPhantomRead(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping api test in short mode")
	}

	t.Cleanup(func() {
		err := os.RemoveAll("./tests/testdata")
		require.NoError(t, err)
	})

	connection, err := net.Dial("tcp", masterAddr)
	require.NoError(t, err)
	defer connection.Close() // nolint

	assert.Equal(t, "OK", sendRequest(t, connection, "begin"))
	assert.Equal(t, "NIL", sendRequest(t, connection, "get random_key"))
	go func() {
		connection, err := net.Dial("tcp", masterAddr)
		require.NoError(t, err)
		defer connection.Close() // nolint

		assert.Equal(t, "OK", sendRequest(t, connection, "begin"))
		assert.Equal(t, "OK", sendRequest(t, connection, "set random_key value"))
		assert.Equal(t, "OK", sendRequest(t, connection, "commit"))
	}()
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, "NIL", sendRequest(t, connection, "get random_key"))
	assert.Equal(t, "OK", sendRequest(t, connection, "commit"))
}

func TestTransactionWriteSkew(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping api test in short mode")
	}

	t.Cleanup(func() {
		err := os.RemoveAll("./tests/testdata")
		require.NoError(t, err)
	})

	connection, err := net.Dial("tcp", masterAddr)
	require.NoError(t, err)
	defer connection.Close() // nolint

	assert.Equal(t, "OK", sendRequest(t, connection, "begin"))
	if assert.Equal(t, "NIL", sendRequest(t, connection, "get key1")) {
		assert.Equal(t, "OK", sendRequest(t, connection, "set key2 val2"))
	}

	go func() {
		connection, err := net.Dial("tcp", masterAddr)
		require.NoError(t, err)
		defer connection.Close() // nolint
		assert.Equal(t, "OK", sendRequest(t, connection, "begin"))
		if assert.Equal(t, "NIL", sendRequest(t, connection, "get key2")) {
			assert.Equal(t, "OK", sendRequest(t, connection, "set key1 val1"))
		}
		assert.Equal(t, "OK", sendRequest(t, connection, "commit"))
	}()
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, "OK", sendRequest(t, connection, "commit"))
	assert.Equal(t, "val1", sendRequest(t, connection, "get key1"))
	assert.Equal(t, "val2", sendRequest(t, connection, "get key2"))
}
