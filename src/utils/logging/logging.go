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

// Initialize the logging objects
func Init(loglevel int) error {
	logLevel = loglevel

	var f *os.File
	if config.LOGGING_FILE != "" {
		var err error
		f, err = os.OpenFile(config.LOGGING_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return err
		}
	} else {
		f = os.Stdout
	}

	errorLogger = log.New(f, "ERROR: ", log.Ldate|log.Ltime)
	warningLogger = log.New(f, "WARNING: ", log.Ldate|log.Ltime)
	infoLogger = log.New(f, "INFO: ", log.Ldate|log.Ltime)
	debugLogger = log.New(f, "DEBUG: ", log.Ldate|log.Ltime)
	fatalLogger = log.New(f, "FATAL: ", log.Ldate|log.Ltime)

	return nil
}

// Logs a message if the log level is bigger than 0
func Fatal(v ...interface{}) {
	if logLevel < LEVEL_FATAL {
		return
	}

	fatalLogger.Fatalln(v...)
}

// Logs a message if the log level is bigger than 0
func Fatalf(format string, v ...interface{}) {
	if logLevel < LEVEL_FATAL {
		return
	}

	fatalLogger.Fatalf(format, v...)
}

// Logs a message if the log level is bigger than 1
func Error(v ...interface{}) {
	if logLevel < LEVEL_ERROR {
		return
	}

	errorLogger.Println(v...)
}

// Logs a message if the log level is bigger than 1
func Errorf(format string, v ...interface{}) {
	if logLevel < LEVEL_ERROR {
		return
	}

	errorLogger.Printf(format, v...)
}

// Logs a message if the log level is bigger than 2
func Warning(v ...interface{}) {
	if logLevel < LEVEL_WARNING {
		return
	}

	warningLogger.Println(v...)
}

// Logs a message if the log level is bigger than 2
func Warningf(format string, v ...interface{}) {
	if logLevel < LEVEL_WARNING {
		return
	}

	warningLogger.Printf(format, v...)
}

// Logs a message if the log level is bigger than 3
func Info(v ...interface{}) {
	if logLevel < LEVEL_INFO {
		return
	}

	infoLogger.Println(v...)
}

// Logs a message if the log level is bigger than 3
func Infof(format string, v ...interface{}) {
	if logLevel < LEVEL_INFO {
		return
	}

	infoLogger.Printf(format, v...)
}

// Logs a message if the log level is bigger than 4
func Debug(v ...interface{}) {
	if logLevel < LEVEL_DEBUG {
		return
	}

	debugLogger.Println(v...)
}

// Logs a message if the log level is bigger than 4
func Debugf(format string, v ...interface{}) {
	if logLevel < LEVEL_DEBUG {
		return
	}

	debugLogger.Printf(format, v...)
}
