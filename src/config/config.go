package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/dattito/purrmannplus-backend/utils"
)

// ENVIRONMENT VARIABLES
var (
	DOT_ENV_FILE_PATH                             string // Path to the .env file, only needed if USE_DOT_ENV_FILE is set to true
	USE_DOT_ENV_FILE                              bool   // If true, the .env file will be loaded
	DATABASE_LOG_LEVEL                            int    // Log level for the database: 1-4: 1:Silent, 2:Error, 3:Warn, 4: Info
	LISTENING_PORT                                int    // The port the api will listen on
	API_URL                                       string // The url the api will be available at, used for the phone number confirmation message link
	CORS_ALLOWED_ORIGINS                          string // Comma separated list of allowed origins or "*"
	AUTHORIZATION_COOKIE_DOMAIN                   string // If set, in the authorization cookie will be set the domain
	AUTHORIZATION_COOKIE_HTTPONLY                 bool   // If true, the cookie will be set as httponly
	AUTHORIZATION_COOKIE_SECURE                   bool   // If true, the cookie will be set as secure
	AUTHORIZATION_COOKIE_SAMESITE                 string // Either "lax" (default), "strict", "disabled" or "none"
	AUTHORIZATION_EXPIRATION_TIME                 int    // The expiration time of the auth jwt tokens / auth cookies in seconds (default to 2678400 seconds = 1 month)
	ENABLE_API                                    bool   // If true, the api will be enabled, otherwise there will be no listener
	ENABLE_SUBSTITUTIONS_SCHEDULER                bool   // If true, the substitutions scheduler will be enabled
	SUBSTITUTIONS_UPDATECRON                      string // Cron expression for the substitutions scheduler
	MAX_ERROS_TO_STOP_UPDATING_SUBSTITUTIONS      int    // If the substitutions scheduler encounters more than this number of errors, it will stop
	MAX_ERROS_TO_STOP_UPDATING_MOODLE_ASSIGNMENTS int    // If the substitutions scheduler encounters more than this number of errors, it will stop
	MOODLE_UPDATECRON                             string // Cron expression for the moodle scheduler
	DATABASE_URI                                  string // The database uri in the format of the given database type
	DATABASE_TYPE                                 string // The database type: SQLITE, POSTGRES, MYSQL
	DATABASE_AUTOMIGRATE                          bool   // If true, the database will be automatically migrated on startup
	SIGNAL_CLI_GRPC_API_URL                       string // The url of the signal cli grpc api
	SIGNAL_SENDER_PHONENUMBER                     string // The phonenumber of the signal sender
	JWT_SECRET                                    string // The secret used to sign the jwt tokens
	SUBSTITUTION_URL                              string // The url of the substitution website
	MOODLE_URL                                    string // The url of the moodle website
	LOGGING_FILE                                  string // The file to log to, if empty, logs to stdout
	LOG_LEVEL                                     int    // 0-5: 0:silent, 1:fatal, 2:error, 3:warn, 4:info, 5:debug
	PATH_TO_API_VIEWS                             string // The path to the api views, default is "./api/providers/rest/views"
	PATH_TO_API_STATIC                            string // The path to the static files of the api, default is "./api/providers/rest/static"
	CONTACT_EMAIL                                 string // The email address users can send emails to
	CONTACT_INSTAGRAM                             string // The instagram account users can send messages to
)

// END OF ENDVIRONMENT VARIABLES

var JWT_SHORTLIVING_SECRET string

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

	CORS_ALLOWED_ORIGINS = utils.GetEnv("CORS_ALLOWED_ORIGINS", "")

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

	AUTHORIZATION_COOKIE_SAMESITE = utils.GetEnv("AUTHORIZATION_COOKIE_SAMESITE", "lax")
	if !utils.Contains([]string{"lax", "strict", "disabled", "none"}, strings.ToLower(AUTHORIZATION_COOKIE_SAMESITE)) {
		return fmt.Errorf("AUTHORIZATION_COOKIE_SAMESITE must be one of lax, strict, disabled, none")
	}

	AUTHORIZATION_EXPIRATION_TIME, err = utils.GetIntEnv("AUTH_EXPIRATION_TIME", 2678400)
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

	MAX_ERROS_TO_STOP_UPDATING_MOODLE_ASSIGNMENTS, err = utils.GetIntEnv("MAX_ERROS_TO_STOP_UPDATING_MOODLE_ASSIGNMENTS", 5)
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

	JWT_SHORTLIVING_SECRET = utils.GenerateString(128)

	SUBSTITUTION_URL = utils.GetEnv("SUBSTITUTION_URL", "")

	MOODLE_URL = utils.GetEnv("MOODLE_URL", "https://moodle.hpg-speyer.de")

	LOGGING_FILE = utils.GetEnv("LOGGING_FILE", "")

	LOG_LEVEL, err = utils.GetIntEnv("LOG_LEVEL", 2)
	if err != nil {
		return err
	}

	PATH_TO_API_VIEWS = utils.GetEnv("PATH_TO_API_VIEWS", "./api/providers/rest/views")

	PATH_TO_API_STATIC = utils.GetEnv("PATH_TO_API_STATIC", "./api/providers/rest/static")

	CONTACT_EMAIL = utils.GetEnv("CONTACT_EMAIL", "")

	CONTACT_INSTAGRAM = utils.GetEnv("CONTACT_INSTAGRAM", "")

	return nil
}
