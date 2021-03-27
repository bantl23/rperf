package run

import (
	"fmt"
	"time"
)

type Client struct {
	Addresses   []string
	Connections uint64
	Duration    time.Duration
	Size        uint64
}

func (client *Client) Run() error {
	fmt.Println(client)
	return nil
}
