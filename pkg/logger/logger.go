package logger

import (
	"log"
	"os"
)

var (
	infoLog  *log.Logger
	errorLog *log.Logger
)

func init() { //nolint:gochecknoinits
	infoLog = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	errorLog = log.New(os.Stderr,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(v ...any) {
	infoLog.Println(v...)
}

func Error(v ...any) {
	errorLog.Println(v...)
}

func Fatal(v ...any) {
	errorLog.Fatal(v...)
}

func Infof(format string, args ...any) {
	infoLog.Printf(format, args...)
}

func Debugf(format string, args ...any) {
	infoLog.Printf(format, args...)
}

func Errorf(format string, args ...any) {
	errorLog.Printf(format, args...)
}

func Fatalf(format string, args ...any) {
	errorLog.Printf(format, args...)
}
