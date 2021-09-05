package logger

import (
	"io"
	"time"
)

var instance *Logger
var inited bool

func Init(path, prefix, delimiter string) {
	instance = NewLogger(path, prefix, delimiter)
	inited = true
}

func GetInstance() *Logger {
	if !inited {
		Init("", "", "")
	}

	return instance
}

func SetAsync(async bool) {
	if !inited {
		Init("", "", "")
	}

	instance.SetAsync(async)
}

func SetStdout(stdout bool) {
	if !inited {
		Init("", "", "")
	}

	instance.SetStdout(stdout)
}

func GetWriter(grade LogGrade) io.Writer {
	if !inited {
		Init("", "", "")
	}
	return instance.GetWriter(grade)
}

func Message(message string, grade LogGrade) {
	if !inited {
		Init("", "", "")
	}

	instance.Message(message, grade)
}

func MessageF(message string, grade LogGrade, args ...interface{}) {
	if !inited {
		Init("", "", "")
	}

	instance.MessageF(message, grade, args...)
}

func Error(err error) {
	if !inited {
		Init("", "", "")
	}

	instance.Error(err)
}

func Notice(message string) {
	if !inited {
		Init("", "", "")
	}

	instance.Notice(message)
}

func NoticeF(format string, args ...interface{}) {
	if !inited {
		Init("", "", "")
	}

	instance.NoticeF(format, args...)
}

func Warning(message string) {
	if !inited {
		Init("", "", "")
	}

	instance.Warning(message)
}

func WarningF(format string, args ...interface{}) {
	if !inited {
		Init("", "", "")
	}

	instance.WarningF(format, args...)
}

func getLogFilename(prefix string) string {
	return prefix + time.Now().Format("02-01-2006") + ".log"
}
