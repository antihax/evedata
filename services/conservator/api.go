package conservator

import (
	"net"
	"net/http"
	"net/rpc"
)

func (s *Conservator) runRPC() error {
	l, err := net.Listen("tcp", "0.0.0.0:3001")
	if err != nil {
		return err
	}

	err = rpc.Register(s)
	if err != nil {
		return err
	}

	go http.Serve(l, nil)
	return nil
}

type Args struct {
	A, B int
}

func (s *Conservator) AddDiscord(args *Args, reply *int) error {
	return nil
}
