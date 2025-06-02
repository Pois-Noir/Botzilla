package botzilla

import (
	"net"
	"sync"
)

type Tunnel struct {
	Name   string
	Open   bool
	source chan byte
	mu     sync.Mutex
	conn   net.Listener
}

func NewTunnel(name string, Source chan byte) *Tunnel {

	new_tunnel := &Tunnel{
		Name:   name,
		source: Source,
		Open:   true,
	}

	go new_tunnel.stream()
	return new_tunnel

}

func (t *Tunnel) ChangeSource(source chan byte) {

	t.Open = false

	// making sure previous connection is closed while changing source
	t.mu.Lock()
	t.source = source
	t.mu.Unlock()

	go t.stream()
}

func (t *Tunnel) stream() {

	// making sure there is only one source of input for
	// tunnel at a time
	t.mu.Lock()
	defer t.mu.Unlock()

	for t.Open {

	}

}
