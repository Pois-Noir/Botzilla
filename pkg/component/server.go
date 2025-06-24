package component

import (
	"github.com/grandcat/zeroconf"
	"sync"
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

func
