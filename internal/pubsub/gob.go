package pubsub

import (
	"bytes"
	"encoding/gob"
)

func encode[T any](val T) ([]byte, error) {
	// Create an encoder to send a value.
	var encodedBytes bytes.Buffer
	enc := gob.NewEncoder(&encodedBytes)
	err := enc.Encode(val)
	if err != nil {
		return nil, err
	}
	return encodedBytes.Bytes(), nil
}

func decode[T any](data []byte) (T, error) {
	dataBuf := bytes.NewBuffer(data)
	var val T
	dec := gob.NewDecoder(dataBuf)
	err := dec.Decode(&val)
	if err != nil {
		return val, err
	}
	return val, nil
}
