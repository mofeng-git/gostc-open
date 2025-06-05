package common

import (
	"fmt"
	"github.com/go-gost/core/logger"
	"github.com/go-gost/x/config/parsing"
	xlogger "github.com/go-gost/x/logger"
	arpcLog "github.com/lesismal/arpc/log"
	frpLog "gostc-sub/pkg/p2p/pkg/util/log"
	"log"
	"sync"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	level := logger.InfoLevel
	//switch logLevel {
	//case "trace":
	//	level = logger.TraceLevel
	//case "debug":
	//	level = logger.DebugLevel
	//case "info":
	//	level = logger.InfoLevel
	//case "warn":
	//	level = logger.WarnLevel
	//case "error":
	//	level = logger.ErrorLevel
	//case "fatal":
	//	level = logger.FatalLevel
	//}
	writer := NewLogWriter()
	logger.SetDefault(xlogger.NewLogger(xlogger.LevelOption(level), xlogger.OutputOption(writer)))
	tlsConfig, _ := parsing.BuildDefaultTLSConfig(nil)
	parsing.SetDefaultTLSConfig(tlsConfig)
	frpLog.RefreshDefault()
	arpcLog.SetLevel(arpcLog.LevelError)
}

var Logger = NewCircularLogger(10000)

type LogEntry struct {
	Timestamp time.Time
	Type      string
	Message   string
}

type CircularLogger struct {
	mu      sync.Mutex
	logs    []LogEntry
	maxSize int
	console bool
}

func NewCircularLogger(maxSize int) *CircularLogger {
	return &CircularLogger{
		logs:    make([]LogEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

func (cl *CircularLogger) AddLog(tp, message string) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now(),
		Message:   message,
		Type:      tp,
	}
	if cl.console {
		//fmt.Println("[LOG]", "Timestamp:", time.Now().Format(time.DateTime))
		//fmt.Println("[LOG]", "Type:", tp)
		//fmt.Println("[LOG]", "Message:", message)
		fmt.Println(fmt.Sprintf("[LOG] %s %s: %s", time.Now().Format(time.DateTime), tp, message))
	}

	if len(cl.logs) < cl.maxSize {
		cl.logs = append(cl.logs, entry)
	} else {
		// 移除最早的日志，添加新的日志
		cl.logs = append(cl.logs[1:], entry)
	}
}

func (cl *CircularLogger) GetLogs() []LogEntry {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	logs := make([]LogEntry, len(cl.logs))
	copy(logs, cl.logs)
	return logs
}

func (cl *CircularLogger) ClearLogs() {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.logs = make([]LogEntry, 0, cl.maxSize)
}

func (cl *CircularLogger) Console(b bool) {
	cl.console = b
}
