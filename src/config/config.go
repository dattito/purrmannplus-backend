package config

import (
	"os"

	"github.com/datti-to/purrmannplus-backend/utils"
)

var DOT_ENV_FILE_PATH string
var USE_DOT_ENV_FILE bool

var LISTENING_PORT int

var SUBSTITUTIONS_UPDATECRON string
var MOODLE_UPDATECRON string

var DATABASE_URI string
var DATABASE_TYPE string
var DATABASE_AUTOMIGRATE bool

var SIGNAL_CLI_GRPC_API_URL string
var SIGNAL_SENDER_PHONENUMBER string

var JWT_SECRET string

var SUBSTITUTION_URL string

func InitConfig() error {

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

	SUBSTITUTIONS_UPDATECRON = utils.GetEnv("SUBSTITUTIONS_UPDATECRON", "*/10 6-23 * * *")
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

	SUBSTITUTION_URL, err = utils.GetEnvInDev("SUBSTITUTION_URL", "https://vertretungsplan.hpg-speyer.de")
	if err != nil {
		return err
	}

	return nil
}
