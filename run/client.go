package run

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Client struct {
	Addresses   []string
	Connections int
	Duration    time.Duration
	Size        int
	Udp         bool
}

func (client *Client) Run() error {
	buffer := make([]byte, client.Size)

	connectChans := make([]chan struct{}, 0)
	statsChan := make([]chan Stats, 0)
	beginChans := make([]chan struct{}, 0)
	endChans := make([]chan struct{}, 0)

	wg := &sync.WaitGroup{}

	for i := 0; i < len(client.Addresses); i++ {
		for j := 0; j < client.Connections; j++ {
			connectChan := make(chan struct{}, 1)
			beginChan := make(chan struct{}, 1)
			statChan := make(chan Stats, 1)
			endChan := make(chan struct{}, 1)
			wg.Add(1)
			go clientRunTcp(client.Addresses[i], j, buffer, connectChan, statChan, beginChan, endChan, wg)
			connectChans = append(connectChans, connectChan)
			beginChans = append(beginChans, beginChan)
			statsChan = append(statsChan, statChan)
			endChans = append(endChans, endChan)
		}
	}

	// wait for all go routines to connect
	for i := range connectChans {
		<-connectChans[i]
	}

	fmt.Printf("CLIENT: ALL CONNECTED\n")

	// start go routines
	for i := range beginChans {
		beginChans[i] <- struct{}{}
	}

	time.Sleep(client.Duration)

	// stop go routines
	for i := range endChans {
		endChans[i] <- struct{}{}
	}

	wg.Wait()

	totals := make(map[string]*Stats)
	items := make(map[string]int)
	// calculate bandwidth
	for i := range statsChan {
		s := <-statsChan[i]
		_, ok := totals[s.Address]
		if !ok {
			items[s.Address] = 0
			totals[s.Address] = &Stats{
				Address:     s.Address,
				Bytes:       s.Bytes,
				ElapsedTime: s.ElapsedTime,
			}
			items[s.Address] = items[s.Address] + 1
			totals[s.Address].Bytes = totals[s.Address].Bytes + s.Bytes
			totals[s.Address].ElapsedTime = totals[s.Address].ElapsedTime + s.ElapsedTime
		}
	}
	for k := range totals {
		mbps := float64(totals[k].Bytes) * 8 / 1024 / 1024 / totals[k].ElapsedTime.Seconds() * float64(items[k])
		fmt.Printf("CLIENT AVG [%s]: %f Mbps\n", totals[k].Address, mbps)
	}

	return nil
}

func clientRunTcp(address string, item int, buffer []byte, connected chan<- struct{}, stats chan<- Stats, beg <-chan struct{}, end <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	fmt.Printf("CLIENT %s: CONNECTED to %s [%+v]\n", conn.LocalAddr(), conn.RemoteAddr(), time.Now())

	connected <- struct{}{} // notify connected
	<-beg                   // block until all routines are connected

	fmt.Printf("CLIENT %s: STARTING with %s [%+v]\n", conn.LocalAddr(), conn.RemoteAddr(), time.Now())

	totalBytes := uint64(0)
	totalElapsed := time.Duration(0)

	done := false
	for !done {
		now := time.Now()
		n, err := conn.Write(buffer)
		if err != nil {
			done = true
			continue
		}
		elapsed := time.Since(now)

		totalBytes = totalBytes + uint64(n)
		totalElapsed = totalElapsed + elapsed

		select {
		case <-end:
			done = true
		default:
		}
	}
	stats <- Stats{
		Address:     conn.RemoteAddr().String(),
		Item:        item,
		Bytes:       totalBytes,
		ElapsedTime: totalElapsed,
	}

	mbps := float64(totalBytes) * 8 / 1024 / 1024 / totalElapsed.Seconds()
	fmt.Printf("CLIENT [%s]: %f Mbps\n", conn.LocalAddr(), mbps)
}
