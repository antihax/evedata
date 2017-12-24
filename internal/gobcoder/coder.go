package gobcoder

import (
	"gopkg.in/mgo.v2/bson"
)

// GobDecoder decode a []byte into a struct
func GobDecoder(message []byte, s interface{}) error {
	return bson.Unmarshal(message, s)
}

// GobEncoder encodes a struct into a []byte
func GobEncoder(s interface{}) ([]byte, error) {
	return bson.Marshal(s)
}
