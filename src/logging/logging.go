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

	logLevel int
)

const (
	LEVEL_DEBUG   = 5
	LEVEL_INFO    = 4
	LEVEL_WARNING = 3
	LEVEL_ERROR   = 2
	LEVEL_FATAL   = 1
	LEVEL_SILENT  = 0
)

func Init(loglevel int) {
	logLevel = loglevel

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

	errorLogger = log.New(f, "ERROR: ", log.Ldate|log.Ltime|log.Llongfile)
	warningLogger = log.New(f, "WARNING: ", log.Ldate|log.Ltime|log.Llongfile)
	infoLogger = log.New(f, "INFO: ", log.Ldate|log.Ltime|log.Llongfile)
	debugLogger = log.New(f, "DEBUG: ", log.Ldate|log.Ltime|log.Llongfile)
	fatalLogger = log.New(f, "FATAL: ", log.Ldate|log.Ltime|log.Llongfile)
}

func Fatal(v ...interface{}) {
	if logLevel < LEVEL_FATAL {
		return
	}

	fatalLogger.Fatalln(v...)
}

func Fatalf(format string, v ...interface{}) {
	if logLevel < LEVEL_FATAL {
		return
	}

	fatalLogger.Fatalf(format, v...)
}

func Error(v ...interface{}) {
	if logLevel < LEVEL_ERROR {
		return
	}

	errorLogger.Println(v...)
}

func Errorf(format string, v ...interface{}) {
	if logLevel < LEVEL_ERROR {
		return
	}

	errorLogger.Printf(format, v...)
}

func Warning(v ...interface{}) {
	if logLevel < LEVEL_WARNING {
		return
	}

	warningLogger.Println(v...)
}

func Warningf(format string, v ...interface{}) {
	if logLevel < LEVEL_WARNING {
		return
	}

	warningLogger.Printf(format, v...)
}

func Info(v ...interface{}) {
	if logLevel < LEVEL_INFO {
		return
	}

	infoLogger.Println(v...)
}

func Infof(format string, v ...interface{}) {
	if logLevel < LEVEL_INFO {
		return
	}

	infoLogger.Printf(format, v...)
}

func Debug(v ...interface{}) {
	if logLevel < LEVEL_DEBUG {
		return
	}

	debugLogger.Println(v...)
}

func Debugf(format string, v ...interface{}) {
	if logLevel < LEVEL_DEBUG {
		return
	}

	debugLogger.Printf(format, v...)
}
