package config

import (
	"fmt"
	"os"

	"github.com/dattito/purrmannplus-backend/utils"
)

var (
	DOT_ENV_FILE_PATH                        string
	USE_DOT_ENV_FILE                         bool
	DATABASE_LOG_LEVEL                       int // 0-5: 0:silent, 1:fatal, 2:error, 3:warn, 4:info, 5:debug
	LISTENING_PORT                           int
	API_URL                                  string
	AUTHORIZATION_COOKIE_DOMAIN              string
	ENABLE_API                               bool
	ENABLE_SUBSTITUTIONS_SCHEDULER           bool
	SUBSTITUTIONS_UPDATECRON                 string
	MAX_ERROS_TO_STOP_UPDATING_SUBSTITUTIONS int
	MOODLE_UPDATECRON                        string
	DATABASE_URI                             string
	DATABASE_TYPE                            string
	DATABASE_AUTOMIGRATE                     bool
	SIGNAL_CLI_GRPC_API_URL                  string
	SIGNAL_SENDER_PHONENUMBER                string
	JWT_SECRET                               string
	JWT_RANDOM_SECRET                        string
	SUBSTITUTION_URL                         string
	LOGGING_FILE                             string
	LOG_LEVEL                                int
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

	JWT_RANDOM_SECRET = utils.GenerateString(128)

	SUBSTITUTION_URL, err = utils.GetEnvInDev("SUBSTITUTION_URL", "https://vertretungsplan.hpg-speyer.de")
	if err != nil {
		return err
	}

	LOGGING_FILE = utils.GetEnv("LOGGING_FILE", "")

	LOG_LEVEL, err = utils.GetIntEnv("LOG_LEVEL", 2)
	if err != nil {
		return err
	}

	return nil
}
