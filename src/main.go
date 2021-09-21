package main

import (
	"log"

	"github.com/datti-to/purrmannplus-backend/api"
	"github.com/datti-to/purrmannplus-backend/app"
	"github.com/datti-to/purrmannplus-backend/config"
	"github.com/datti-to/purrmannplus-backend/database"
	"github.com/datti-to/purrmannplus-backend/services/scheduler"
	"github.com/datti-to/purrmannplus-backend/services/signal_message_sender"
)

func main() {
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

	//scheduler.StartScheduler()

	log.Fatal(api.StartListening())
}
