package logger

type LoggerWriter struct {
	l     *FileLogger
	grade LogGrade
}

func NewLoggerWriter(l *FileLogger, grade LogGrade) *LoggerWriter {
	return &LoggerWriter{l, grade}
}

func (w *LoggerWriter) Write(p []byte) (n int, err error) {
	w.l.Message(string(p), w.grade)
	return len(p), nil
}
