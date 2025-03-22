package wal

import (
	"bytes"
	"encoding/gob"
)

type Command struct {
	LSN         int
	CommandType int
	Args        []string
}

func (c *Command) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	if err := encoder.Encode(*c); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Command) Decode(buf *bytes.Buffer) error {
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(c); err != nil {
		return err
	}

	return nil
}
