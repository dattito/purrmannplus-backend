package logging

import (
	"log"
	"os"

	"github.com/dattito/purrmannplus-backend/config"
)

var (
	errorLogger   *log.Logger
	warningLogger *log.Logger
	infoLogger    *log.Logger
	debugLogger   *log.Logger
	fatalLogger   *log.Logger
)

func Init() {
	var f *os.File
	if config.LOGGING_FILE != "" {
		var err error
		f, err = os.OpenFile(config.LOGGING_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
	} else {
		f = os.Stdout
	}

	errorLogger = log.New(f, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLogger = log.New(f, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(f, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger = log.New(f, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	fatalLogger = log.New(f, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Error(v ...interface{}) {
	errorLogger.Println(v...)
}

func Errorf(format string, v ...interface{}) {
	errorLogger.Printf(format, v...)
}

func Warning(v ...interface{}) {
	warningLogger.Println(v...)
}

func Warningf(format string, v ...interface{}) {
	warningLogger.Printf(format, v...)
}

func Info(v ...interface{}) {
	infoLogger.Println(v...)
}

func Infof(format string, v ...interface{}) {
	infoLogger.Printf(format, v...)
}

func Debug(v ...interface{}) {
	debugLogger.Println(v...)
}

func Debugf(format string, v ...interface{}) {
	debugLogger.Printf(format, v...)
}

func Fatal(v ...interface{}) {
	fatalLogger.Fatalln(v...)
}

func Fatalf(format string, v ...interface{}) {
	fatalLogger.Fatalf(format, v...)
}
