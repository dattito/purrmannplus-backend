package config

import "github.com/datti-to/purrmannplus-backend/utils"

var SUBSTITUTIONS_UPDATECRON string = utils.GetEnv("SUBSTITUTIONS_UPDATECRON", "*/10 6-23 * * *")
var MOODLE_UPDATECRON string = utils.GetEnv("MOODLE_UPDATECRON", "0 6-23 * * *")

var DATABASE_URI string = utils.GetEnv("DATABASE_URI", "db.sqlite")
var DATABASE_TYPE string = utils.GetEnv("DATABASE_TYPE", "SQLITE")
var SIGNAL_CLI_GRPC_API_URL string = utils.GetEnvInDev("SIGNAL_CLI_GRPC_API_URL", "localhost:9000")
var DATABASE_AUTOMIGRATE bool = utils.GetBoolEnv("DATABASE_AUTOMIGRATE", true)
