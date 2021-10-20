package config

import (
	"fmt"
	"os"

	"github.com/dattito/purrmannplus-backend/utils"
)

var (
	DOT_ENV_FILE_PATH                        string // Path to the .env file, only needed if USE_DOT_ENV_FILE is set to true
	USE_DOT_ENV_FILE                         bool   // If true, the .env file will be loaded
	DATABASE_LOG_LEVEL                       int    // Log level for the database: 1-4: 1:Silent, 2:Error, 3:Warn, 4: Info
	LISTENING_PORT                           int    // The port the api will listen on
	API_URL                                  string // The url the api will be available at, used for the phone number confirmation message link
	AUTHORIZATION_COOKIE_DOMAIN              string // If set, in the authorization cookie will be set the domain
	AUTHORIZATION_COOKIE_HTTPONLY            bool   // If true, the cookie will be set as httponly
	AUTHORIZATION_COOKIE_SECURE              bool   // If true, the cookie will be set as secure
	ENABLE_API                               bool   // If true, the api will be enabled, otherwise there will be no listener
	ENABLE_SUBSTITUTIONS_SCHEDULER           bool   // If true, the substitutions scheduler will be enabled
	SUBSTITUTIONS_UPDATECRON                 string // Cron expression for the substitutions scheduler
	MAX_ERROS_TO_STOP_UPDATING_SUBSTITUTIONS int    // If the substitutions scheduler encounters more than this number of errors, it will stop
	MOODLE_UPDATECRON                        string // Cron expression for the moodle scheduler
	DATABASE_URI                             string // The database uri in the format of the given database type
	DATABASE_TYPE                            string // The database type: SQLITE, POSTGRES, MYSQL
	DATABASE_AUTOMIGRATE                     bool   // If true, the database will be automatically migrated on startup
	SIGNAL_CLI_GRPC_API_URL                  string // The url of the signal cli grpc api
	SIGNAL_SENDER_PHONENUMBER                string // The phonenumber of the signal sender
	JWT_SECRET                               string // The secret used to sign the jwt tokens
	DNT_JWT_RANDOM_SECRET                    string // Gets generated on init, so DO NOT TOUCH
	SUBSTITUTION_URL                         string // The url of the substitution website
	LOGGING_FILE                             string // The file to log to, if empty, logs to stdout
	LOG_LEVEL                                int    // 0-5: 0:silent, 1:fatal, 2:error, 3:warn, 4:info, 5:debug
	DNT_VERSION                              string // The version of this application. It is automatically set by docker, so DO NOT TOUCH
)

func Init() error {

	var err error

	DOT_ENV_FILE_PATH = utils.GetEnv("DOT_ENV_FILE_PATH", ".env")

	USE_DOT_ENV_FILE, err = utils.GetBoolEnv("USE_DOT_ENV_FILE", true)
	if err != nil {
		return err
	}

	if USE_DOT_ENV_FILE {
		if _, err := os.Stat(DOT_ENV_FILE_PATH); !os.IsNotExist(err) {
			err = utils.LoadDotEnvFile()
			if err != nil {
				return err
			}
		}
	}

	DATABASE_LOG_LEVEL, err = utils.GetIntEnv("DATABASE_LOG_LEVEL", 1)
	if err != nil {
		return err
	}

	LISTENING_PORT, err = utils.GetIntEnv("LISTENING_PORT", 3000)
	if err != nil {
		return err
	}

	API_URL = utils.GetEnv("API_URL", fmt.Sprintf("http://localhost:%d", LISTENING_PORT))

	// If set, in the authorization cookie will be set the domain
	AUTHORIZATION_COOKIE_DOMAIN = utils.GetEnv("AUTHORIZATION_COOKIE_DOMAIN", "")

	AUTHORIZATION_COOKIE_HTTPONLY, err = utils.GetBoolEnv("AUTHORIZATION_COOKIE_HTTPONLY", false)
	if err != nil {
		return err
	}

	AUTHORIZATION_COOKIE_SECURE, err = utils.GetBoolEnv("AUTHORIZATION_COOKIE_SECURE", false)
	if err != nil {
		return err
	}

	ENABLE_API, err = utils.GetBoolEnv("ENABLE_API", true)
	if err != nil {
		return err
	}

	ENABLE_SUBSTITUTIONS_SCHEDULER, err = utils.GetBoolEnv("ENABLE_SUBSTITUTIONS_SCHEDULER", true)
	if err != nil {
		return err
	}
	SUBSTITUTIONS_UPDATECRON = utils.GetEnv("SUBSTITUTIONS_UPDATECRON", "*/10 6-23 * * *")
	MAX_ERROS_TO_STOP_UPDATING_SUBSTITUTIONS, err = utils.GetIntEnv("MAX_ERROS_TO_STOP_UPDATING_SUBSTITUTIONS", 5)
	if err != nil {
		return err
	}
	MOODLE_UPDATECRON = utils.GetEnv("MOODLE_UPDATECRON", "0 6-23 * * *")
	DATABASE_URI = utils.GetEnv("DATABASE_URI", "db.sqlite")
	DATABASE_TYPE = utils.GetEnv("DATABASE_TYPE", "SQLITE")

	DATABASE_AUTOMIGRATE, err = utils.GetBoolEnv("DATABASE_AUTOMIGRATE", true)
	if err != nil {
		return err
	}

	SIGNAL_CLI_GRPC_API_URL, err = utils.GetEnvInDev("SIGNAL_CLI_GRPC_API_URL", "localhost:9000")
	if err != nil {
		return err
	}

	SIGNAL_SENDER_PHONENUMBER, err = utils.GetEnvInDev("SIGNAL_SENDER_PHONENUMBER", "+1555123456")
	if err != nil {
		return err
	}

	JWT_SECRET, err = utils.GetEnvInDev("JWT_SECRET", "secret")
	if err != nil {
		return err
	}

	DNT_JWT_RANDOM_SECRET = utils.GenerateString(128)

	SUBSTITUTION_URL, err = utils.GetEnvInDev("SUBSTITUTION_URL", "https://vertretungsplan.hpg-speyer.de")
	if err != nil {
		return err
	}

	LOGGING_FILE = utils.GetEnv("LOGGING_FILE", "")

	LOG_LEVEL, err = utils.GetIntEnv("LOG_LEVEL", 2)
	if err != nil {
		return err
	}

	// Gots passed by the build script (DNT=DO NOT TOUCH)
	DNT_VERSION = utils.GetEnv("DNT_VERSION", "")

	return nil
}
