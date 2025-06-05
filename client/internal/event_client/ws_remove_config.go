package event_client

import (
	"github.com/go-gost/x/registry"
	registry2 "gostc-sub/pkg/p2p/registry"
)

func (e *Event) WsRemoveConfig(names []string) {
	for _, name := range names {
		if svc := registry.ServiceRegistry().Get(name); svc != nil {
			_ = svc.Close()
			registry.ServiceRegistry().Unregister(name)
		}
		registry2.Del(name)
		e.svcMap.Store(name, true)
	}
}
