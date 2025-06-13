package component

import (
	"fmt"
	"net"
	"sync"

	safemap "github.com/Pois-Noir/Botzilla-Utils/safemap"
)

type tunnel struct {
	Name     string
	source   chan byte
	channels *safemap.SafeMap[net.Addr, chan byte]
	mu       sync.Mutex
	port     int
	Stop     bool
}

func NewTunnel(name string, Source chan byte, p int) *tunnel {

	new_tunnel := &tunnel{
		Name:     name,
		source:   Source,
		port:     p,
		Stop:     false,
		channels: safemap.NewSafeMap[net.Addr, chan byte](),
	}

	return new_tunnel
}

func (t *tunnel) manageSource() {

	for !t.Stop {
		val := <-t.source
		t.channels.ForEach(func(_ net.Addr, c *chan byte) {
			*c <- val
		})
	}

}

func (t *tunnel) Start() error {

	// Make sure tunnel only runs once
	t.mu.Lock()
	defer t.mu.Unlock()

	go t.manageSource()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", t.port))
	if err != nil {
		return err
	}
	defer listener.Close()

	for !t.Stop {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: \n", err)
			continue
		}

		new_chan := make(chan byte, 1024)
		t.channels.Add(conn.RemoteAddr(), &new_chan)

		go t.sendStream(conn, new_chan)
	}

	return nil
}

func (t *tunnel) sendStream(conn net.Conn, data chan byte) error {

	defer conn.Close()
	addr := conn.RemoteAddr()

	for !t.Stop {
		_, err := conn.Write([]byte{<-data})
		if err != nil {
			// Remove Channel
			t.channels.Remove(addr)
			return err
		}
	}

	return nil
}
