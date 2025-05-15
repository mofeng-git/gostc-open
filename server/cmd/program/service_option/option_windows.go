//go:build windows

package service_option

import (
	"github.com/kardianos/service"
)

func MakeOptions() service.KeyValue {
	return service.KeyValue{
		"OnFailureDelayDuration": "10",
	}
}
