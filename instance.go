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

type Logger struct {
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

func NewLogger(path, prefix, delimiter string) *Logger {
	if delimiter == "" {
		delimiter = "\n------||------\n"
	}

	if path == "" {
		path = "logs"
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	l := Logger{
		logpath:   path,
		prefix:    prefix,
		delimiter: delimiter,
	}

	l.updateFile()

	if instance != nil {
		l.SetAsync(instance.async)
		l.SetStdout(instance.stdout)
	}

	return &l
}

func (l *Logger) Copy() *Logger {
	copy := *l
	return &copy
}

func (l *Logger) SetStdout(stdout bool) {
	l.stdout = stdout
}

func (l *Logger) SetAsync(async bool) {
	l.async = async
}

func (l *Logger) SetEntryPrefix(prefix string) {
	l.entryPrefix = "[" + prefix + "] "
}

func (l *Logger) GetLogger() *log.Logger {
	return l.fileLog
}

func (l *Logger) GetWriter(grade LogGrade) io.Writer {
	return NewLoggerWriter(l, grade)
}

func (l *Logger) Message(message string, grade LogGrade) {
	if l.async {
		go l.message(message, grade)
	} else {
		l.message(message, grade)
	}
}

func (l *Logger) MessageF(message string, grade LogGrade, args ...interface{}) {
	message = fmt.Sprintf(message, args...)
	l.Message(message, grade)
}

func (l *Logger) Error(err error) {
	sentry.CaptureException(err)
	l.MessageF("%s: %s\n%s", LG_Error, err, debug.Stack())
}

func (l *Logger) Notice(message string) {
	l.Message(message, LG_Notice)
}

func (l *Logger) NoticeF(format string, args ...interface{}) {
	l.MessageF(format, LG_Notice, args...)
}

func (l *Logger) Warning(message string) {
	l.Message(message, LG_Warning)
}

func (l *Logger) WarningF(format string, args ...interface{}) {
	l.MessageF(format, LG_Warning, args...)
}

func (l *Logger) message(message string, grade LogGrade) {
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

func (l *Logger) updateFile() {
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
