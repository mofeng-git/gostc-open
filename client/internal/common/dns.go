package common

import (
	"context"
	"net"
	"time"
)

func init() {
	fixDNSResolver()
}

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
