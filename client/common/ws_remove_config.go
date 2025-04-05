package common

import (
	"github.com/go-gost/x/registry"
	registry2 "gostc-sub/p2p/registry"
)

func WsRemoveConfig(names []string) {
	for _, name := range names {
		if svc := registry.ServiceRegistry().Get(name); svc != nil {
			_ = svc.Close()
			registry.ServiceRegistry().Unregister(name)
		}
		registry2.Del(name)
		SvcMap[name] = false
	}
}
