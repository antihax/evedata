package conservator

import (
	"errors"
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
	*reply = false
	g, err := s.discord.Guild(*args)
	if err != nil {
		log.Println(err)
		return nil
	}
	if g.Name != "" {
		*reply = true
	}
	return nil
}

func (s *Conservator) VerifyDiscordChannel(args []string, reply *bool) error {
	*reply = false
	g, err := s.discord.Guild(args[0])
	if err != nil {
		return err
	}
	if g.Name != "" {
		return errors.New("Invalid server")
	}
	for _, c := range g.Channels {
		if c.ID == args[1] {
			*reply = true
		}
	}

	return nil
}
