package run

import (
	"fmt"
	"net"
	"time"
)

type Server struct {
	Address string
	Size    int
}

func (server *Server) Run() error {
	l, err := net.Listen("tcp", server.Address)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer l.Close()

	buffer := make([]byte, server.Size)
	for {
		client, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("SERVER %s: CONNECTED from %s\n", client.LocalAddr(), client.RemoteAddr())
		go serverRunTcp(client, buffer)
	}
}

func serverRunTcp(client net.Conn, buffer []byte) {
	totalBytes := float64(0)
	totalElapsed := time.Duration(0)

	done := false
	for !done {
		now := time.Now()
		n, err := client.Read(buffer)
		if err != nil {
			done = true
			continue
		}
		elapsed := time.Since(now)

		totalBytes = totalBytes + float64(n)
		totalElapsed = totalElapsed + elapsed
	}
	mbps := float64(totalBytes) * 8 / 1024 / 1024 / totalElapsed.Seconds()
	fmt.Printf("SERVER [%s]: %f\n", client.RemoteAddr(), mbps)
}
