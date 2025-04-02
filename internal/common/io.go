package common

import (
	"encoding/binary"
	"io"
)

func Read(reader io.Reader, data []byte) (int, error) {
	var l uint32
	if err := binary.Read(reader, binary.LittleEndian, &l); err != nil {
		return 0, err
	}

	dataLen := int(l)
	if dataLen > len(data) {
		return 0, io.ErrShortBuffer
	}

	data = data[:dataLen]
	if _, err := io.ReadFull(reader, data); err != nil {
		return 0, err
	}

	return dataLen, nil
}

func Write(writer io.Writer, data []byte) (int, error) {
	l := uint32(len(data))
	if err := binary.Write(writer, binary.LittleEndian, l); err != nil {
		return 0, err
	}

	return writer.Write(data)
}
