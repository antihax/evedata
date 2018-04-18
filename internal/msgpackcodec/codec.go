package msgpackcodec

import msgpack "gopkg.in/vmihailenco/msgpack.v2"

// MsgPackCodec is a MessagePack codec for gRPC
type MsgPackCodec struct {
}

func (c *MsgPackCodec) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (c *MsgPackCodec) Unmarshal(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}

func (c *MsgPackCodec) String() string {
	return "MsgPackCodec"
}
