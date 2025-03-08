package main

import (
	"context"
	"errors"
	"net"
	"strings"
	"time"
)

func fixDNSResolver() {
	if net.DefaultResolver != nil {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err := net.DefaultResolver.LookupHost(timeoutCtx, "google.com")
		if err == nil {
			return
		}
	}
	net.DefaultResolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if addr == "127.0.0.1:53" || addr == "[::1]:53" {
				addr = "8.8.8.8:53"
			}
			var d net.Dialer
			return d.DialContext(ctx, network, addr)
		},
	}
}

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
