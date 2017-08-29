package gobcoder

import (
	"bytes"
	"encoding/gob"
)

// GobDecoder decode a []byte into a struct
func GobDecoder(message []byte, s interface{}) error {
	b := bytes.NewBuffer(message)
	dec := gob.NewDecoder(b)
	if err := dec.Decode(s); err != nil {
		return err
	}
	return nil
}

// GobEncoder encodes a struct into a []byte
func GobEncoder(s interface{}) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(s)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
