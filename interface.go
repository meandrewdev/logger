package logger

import (
	"io"
	"log"
)

type Logger interface {
	Copy() Logger

	SetStdout(stdout bool)
	SetAsync(async bool)
	SetEntryPrefix(prefix string)

	IsStdout() bool
	IsAsync() bool
	GetEntryPrefix() string

	GetLogger() *log.Logger
	GetWriter(grade LogGrade) io.Writer

	Message(message string, grade LogGrade)
	MessageF(message string, grade LogGrade, args ...interface{})

	Notice(message string)
	NoticeF(format string, args ...interface{})

	Warning(message string)
	WarningF(format string, args ...interface{})

	Error(err error)
	ErrorF(msg string, args ...interface{})
}
