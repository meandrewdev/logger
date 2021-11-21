package logger

import (
	"io"
	"time"
)

var instance Logger
var inited bool

func Init(path, prefix, delimiter string) {
	if inited {
		return
	}

	instance = NewLogger(path, prefix, delimiter)
	inited = true
}

func Reinit(path, prefix, delimiter string) {
	inited = false

	Init(path, prefix, delimiter)
}

func GetInstance() Logger {
	Init("", "", "")

	return instance
}

func SetAsync(async bool) {
	Init("", "", "")

	instance.SetAsync(async)
}

func SetStdout(stdout bool) {
	Init("", "", "")

	instance.SetStdout(stdout)
}

func SetEntryPrefix(prefix string) {
	Init("", "", "")

	instance.SetEntryPrefix(prefix)
}

func GetWriter(grade LogGrade) io.Writer {
	Init("", "", "")
	return instance.GetWriter(grade)
}

func Message(message string, grade LogGrade) {
	Init("", "", "")

	instance.Message(message, grade)
}

func MessageF(message string, grade LogGrade, args ...interface{}) {
	Init("", "", "")

	instance.MessageF(message, grade, args...)
}

func Error(err error) {
	Init("", "", "")

	instance.Error(err)
}

func ErrorF(msg string, args ...interface{}) {
	Init("", "", "")

	instance.ErrorF(msg, args...)
}

func Notice(message string) {
	Init("", "", "")

	instance.Notice(message)
}

func NoticeF(format string, args ...interface{}) {
	Init("", "", "")

	instance.NoticeF(format, args...)
}

func Warning(message string) {
	Init("", "", "")

	instance.Warning(message)
}

func WarningF(format string, args ...interface{}) {
	Init("", "", "")

	instance.WarningF(format, args...)
}

func getLogFilename(prefix string) string {
	return prefix + time.Now().Format("02-01-2006") + ".log"
}
