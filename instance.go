package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	"sync"

	"github.com/getsentry/sentry-go"
)

type FileLogger struct {
	logpath     string
	filename    string
	file        *os.File
	fileLog     *log.Logger
	prefix      string
	delimiter   string
	async       bool
	stdout      bool
	entryPrefix string

	*sync.Mutex
}

func NewLogger(path, prefix, delimiter string) *FileLogger {
	if delimiter == "" {
		delimiter = "\n------||------\n"
	}

	if path == "" {
		path = "logs"
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	l := FileLogger{
		logpath:   path,
		prefix:    prefix,
		delimiter: delimiter,
	}

	l.updateFile()

	if instance != nil {
		l.SetAsync(instance.IsAsync())
		l.SetStdout(instance.IsStdout())
	}

	return &l
}

func (l *FileLogger) Copy() Logger {
	copy := *l
	return &copy
}

func (l *FileLogger) SetStdout(stdout bool) {
	l.stdout = stdout
}

func (l *FileLogger) IsStdout() bool {
	return l.stdout
}

func (l *FileLogger) SetAsync(async bool) {
	l.async = async
}

func (l *FileLogger) IsAsync() bool {
	return l.async
}

func (l *FileLogger) SetEntryPrefix(prefix string) {
	l.entryPrefix = "[" + prefix + "] "
}

func (l *FileLogger) GetEntryPrefix() string {
	return l.entryPrefix
}

func (l *FileLogger) GetLogger() *log.Logger {
	return l.fileLog
}

func (l *FileLogger) GetWriter(grade LogGrade) io.Writer {
	return NewLoggerWriter(l, grade)
}

func (l *FileLogger) Message(message string, grade LogGrade) {
	if l.async {
		go l.message(message, grade)
	} else {
		l.message(message, grade)
	}
}

func (l *FileLogger) MessageF(message string, grade LogGrade, args ...interface{}) {
	message = fmt.Sprintf(message, args...)
	l.Message(message, grade)
}

func (l *FileLogger) Error(err error) {
	sentry.CaptureException(err)
	l.MessageF("%s: %s\n%s", LG_Error, err, debug.Stack())
}

func (l *FileLogger) Notice(message string) {
	l.Message(message, LG_Notice)
}

func (l *FileLogger) NoticeF(format string, args ...interface{}) {
	l.MessageF(format, LG_Notice, args...)
}

func (l *FileLogger) Warning(message string) {
	l.Message(message, LG_Warning)
}

func (l *FileLogger) WarningF(format string, args ...interface{}) {
	l.MessageF(format, LG_Warning, args...)
}

func (l *FileLogger) message(message string, grade LogGrade) {
	l.Lock()
	defer l.Unlock()

	if l.filename != getLogFilename(l.prefix) {
		l.updateFile()
	}

	msg := fmt.Sprintf("%s: %s%s%s", l.entryPrefix, grade, message, l.delimiter)

	l.fileLog.Print(msg)
	if l.stdout {
		log.Print(msg)
	}
}

func (l *FileLogger) updateFile() {
	newName := getLogFilename(l.prefix)

	if l.filename == newName {
		return
	}

	l.filename = newName

	var err error
	l.file, err = os.OpenFile(l.logpath+"/"+l.filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)

	if err != nil {
		panic(err)
	}

	l.fileLog = log.New(l.file, "", log.LstdFlags|log.Lshortfile)
}
