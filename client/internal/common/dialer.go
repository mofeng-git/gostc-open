package common

import (
	"context"
	"errors"
	"net"
	"strings"
)

type MyDialer struct {
	net.Resolver
	ctx context.Context
}

func NewDialer() *MyDialer {
	return &MyDialer{
		Resolver: net.Resolver{
			PreferGo:     true,
			StrictErrors: false,
			Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
				var split = strings.Split(addr, ":")
				if len(split) != 2 {
					return nil, errors.New("addr analysis fail")
				}
				host, err := net.LookupHost(split[0])
				if len(host) == 0 {
					return nil, err
				}
				var d net.Dialer
				return d.DialContext(ctx, network, host[0]+":"+split[1])
			},
		},
		ctx: context.Background(),
	}
}

func (m *MyDialer) Dial(network, addr string) (c net.Conn, err error) {
	return m.Resolver.Dial(m.ctx, network, addr)
}
