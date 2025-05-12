//go:build !linux && !windows

package service_option

import (
	"fmt"
	"github.com/kardianos/service"
)

func MakeOptions() service.KeyValue {
	fmt.Println("is other")
	return service.KeyValue{
		"OnFailureDelayDuration": "10",
	}
}
