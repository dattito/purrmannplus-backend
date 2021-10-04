package main

import (
	"log"

	"github.com/dattito/purrmannplus-backend/api"
	"github.com/dattito/purrmannplus-backend/app"
	"github.com/dattito/purrmannplus-backend/config"
	"github.com/dattito/purrmannplus-backend/database"
	"github.com/dattito/purrmannplus-backend/logging"
	"github.com/dattito/purrmannplus-backend/services/scheduler"
	"github.com/dattito/purrmannplus-backend/services/signal_message_sender"
)

func main() {

	// Load configuration
	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	logging.Init(config.LOG_LEVEL)

	scheduler.Init()

	if err := signal_message_sender.Init(); err != nil {
		log.Fatal(err)
	}

	if err := database.Init(); err != nil {
		log.Fatal(err)
	}

	app.Init()

	api.Init()

	// Start scheduler
	scheduler.StartScheduler()

	// Start API
	logging.Fatal(api.StartListening())
}
