package run

import (
	"fmt"
)

type Server struct {
	Addresses []string
	Size      uint64
}

func (server *Server) Run() error {
	fmt.Println(server)
	return nil
}
