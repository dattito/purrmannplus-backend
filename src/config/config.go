package config

import "github.com/datti-to/purrmannplus-backend/utils"

var SUBSTITUTIONS_UPDATECRON string
var MOODLE_UPDATECRON string

var DATABASE_URI string
var DATABASE_TYPE string
var SIGNAL_CLI_GRPC_API_URL string
var DATABASE_AUTOMIGRATE bool

var JWT_SECRET string

var SUBSTITUTION_URL string = "https://vertretungsplan.hpg-speyer.de"

func InitConfig() error {
	var err error
	SUBSTITUTIONS_UPDATECRON = utils.GetEnv("SUBSTITUTIONS_UPDATECRON", "*/10 6-23 * * *")
	MOODLE_UPDATECRON = utils.GetEnv("MOODLE_UPDATECRON", "0 6-23 * * *")
	DATABASE_URI = utils.GetEnv("DATABASE_URI", "db.sqlite")
	DATABASE_TYPE = utils.GetEnv("DATABASE_TYPE", "SQLITE")
	SIGNAL_CLI_GRPC_API_URL, err = utils.GetEnvInDev("SIGNAL_CLI_GRPC_API_URL", "localhost:9000")
	if err != nil {
		return err
	}

	DATABASE_AUTOMIGRATE, err = utils.GetBoolEnv("DATABASE_AUTOMIGRATE", true)
	if err != nil {
		return err
	}

	JWT_SECRET, err = utils.GetEnvInDev("JWT_SECRET", "secret")
	if err != nil {
		return err
	}

	return nil
}
