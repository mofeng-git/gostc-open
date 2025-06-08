package signal

import (
	"os"
	"os/signal"
	"syscall"
)

func Free() <-chan os.Signal {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	return signalChan
}
