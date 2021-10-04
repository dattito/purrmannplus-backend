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

	logging.Init()

	// Load configuration
	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

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
	log.Fatal(api.StartListening())
}
