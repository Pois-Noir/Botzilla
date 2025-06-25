package botzilla

import (
	"github.com/grandcat/zeroconf"
	"golang.org/x/net/context"
	"strconv"
	"sync"
	"time"
)

var (
	resolver *zeroconf.Resolver
	once     sync.Once
)

func getResolver() *zeroconf.Resolver {
	once.Do(func() {
		r, err := zeroconf.NewResolver(nil)
		resolver = r

		// TODO Check me pls
		if err != nil {
			panic(err)
		}
	})

	return resolver
}

func GetComponents() (map[string]string, error) {

	resolver := getResolver()
	entries := make(chan *zeroconf.ServiceEntry)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	defer cancel()

	err := resolver.Browse(ctx, "_botzilla._tcp", "local.", entries)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)

	for {

		select {
		case entry := <-entries:
			if len(entry.AddrIPv4) == 0 {
				continue // skip this entry
			}
			result[entry.Instance] = entry.AddrIPv4[0].String() + ":" + strconv.Itoa(entry.Port)
		case <-ctx.Done():
			return result, nil
		}
	}
}

func GetComponent(name string) (string, error) {

	resolver := getResolver()
	entries := make(chan *zeroconf.ServiceEntry)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := resolver.Lookup(ctx, name, "_botzilla._tcp", "local.", entries)
	if err != nil {
		return "", err
	}

	select {
	case entry := <-entries:
		return entry.AddrIPv4[0].String() + ":" + strconv.Itoa(entry.Port), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
