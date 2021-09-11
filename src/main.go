package main

import (
	"log"

	"github.com/datti-to/purrmannplus-backend/api"
	"github.com/datti-to/purrmannplus-backend/config"
	"github.com/datti-to/purrmannplus-backend/database"
	"github.com/datti-to/purrmannplus-backend/services/signal_message_sender"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatal(err)
	}

	if err := signal_message_sender.Init(); err != nil {
		log.Fatal(err)
	}

	if err := database.Init(); err != nil {
		log.Fatal(err)
	}

	api.Init()

	log.Fatal(api.StartListening())
}
