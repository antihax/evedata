package conservator

import (
	"log"
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
	rpc.HandleHTTP()
	go http.Serve(l, nil)
	return nil
}

func (s *Conservator) VerifyDiscord(args *string, reply *bool) error {
	g, err := s.discord.Guild(*args)
	if err != nil {
		log.Println(err)
		*reply = false
		return nil
	}
	if g.Name != "" {
		*reply = true
	}
	return nil
}
