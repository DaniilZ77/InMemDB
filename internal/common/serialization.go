package common

import (
	"bytes"
	"encoding/gob"
	"io"
)

func Encode[Object any](object Object) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(object)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func Decode[Object any](buffer *bytes.Buffer) (Object, error) {
	decoder := gob.NewDecoder(buffer)
	var object Object
	err := decoder.Decode(&object)
	if err != nil {
		return object, err
	}

	return object, nil
}

func DecodeOne[Object any](data []byte) (Object, error) {
	return Decode[Object](bytes.NewBuffer(data))
}

func DecodeMany[Array []Object, Object any](data []byte) (Array, error) {
	buffer := bytes.NewBuffer(data)
	var objects Array
	for {
		object, err := Decode[Array](buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}

		objects = append(objects, object...)
	}

	return objects, nil
}
