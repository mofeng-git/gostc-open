package common

import "fmt"

type LogWriter struct {
	console bool
	logChan chan []byte
}

func NewLogWriter(console bool) *LogWriter {
	return &LogWriter{
		console: console,
		logChan: make(chan []byte, 10000),
	}
}

func (l *LogWriter) Write(p []byte) (n int, err error) {
	if l.console {
		fmt.Print(string(p))
	}
	l.logChan <- p
	return len(p), nil
}

func (l *LogWriter) C() <-chan []byte {
	return l.logChan
}
