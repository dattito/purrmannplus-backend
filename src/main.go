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
		log.Fatalf("Failed to load configuration: %v", err)
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

	// Start scheduler
	scheduler.StartScheduler()

	// Start API
	logging.Fatal(api.StartListening())
}
