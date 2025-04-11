package common

type LogWriter struct {
}

func NewLogWriter() *LogWriter {
	return &LogWriter{}
}

func (l *LogWriter) Write(p []byte) (n int, err error) {
	Logger.AddLog("detail", string(p))
	return len(p), nil
}
