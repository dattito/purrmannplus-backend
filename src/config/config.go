package config

import (
	"fmt"
	"os"

	"github.com/datti-to/purrmannplus-backend/utils"
)

var (
	DOT_ENV_FILE_PATH                        string
	USE_DOT_ENV_FILE                         bool
	LISTENING_PORT                           int
	API_URL                                  string
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

	LISTENING_PORT, err = utils.GetIntEnv("LISTENING_PORT", 3000)
	if err != nil {
		return err
	}

	API_URL = utils.GetEnv("API_URL", fmt.Sprintf("http://localhost:%d", LISTENING_PORT))

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

	return nil
}
