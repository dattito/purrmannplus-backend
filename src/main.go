package main

import (
	"log"

	"github.com/dattito/purrmannplus-backend/api"
	"github.com/dattito/purrmannplus-backend/app"
	"github.com/dattito/purrmannplus-backend/config"
	"github.com/dattito/purrmannplus-backend/database"
	"github.com/dattito/purrmannplus-backend/services/scheduler"
	"github.com/dattito/purrmannplus-backend/services/signal_message_sender"
	"github.com/dattito/purrmannplus-backend/utils/logging"
)

func main() {

	// Load configuration
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	if !config.ENABLE_API && !config.ENABLE_SUBSTITUTIONS_SCHEDULER {
		log.Fatal("No API or scheduler enabled. Exiting.")
	}

	if err := logging.Init(); err != nil {
		log.Fatalf("Failed to initialize logging: %s", err)
	}

	scheduler.Init()

	if err := signal_message_sender.Init(); err != nil {
		log.Fatalf("Failed to initialize signal message sender: %s", err)
	}

	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %s", err)
	}

	app.Init()

	api.Init()

	logging.Infof("Starting PurrmannPlus-Backend %s", config.DNT_VERSION)

	if config.ENABLE_API && config.ENABLE_SUBSTITUTIONS_SCHEDULER {
		// Start scheduler
		scheduler.StartAsync()
	}

	if config.ENABLE_API {
		// Start API
		logging.Fatal(api.StartListening())
	}

	if !config.ENABLE_API && config.ENABLE_SUBSTITUTIONS_SCHEDULER {
		// Start scheduler
		scheduler.StartBlocking()
	}
}
