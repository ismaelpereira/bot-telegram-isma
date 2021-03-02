package main

import (
	"log"

	"github.com/ismaelpereira/telegram-bot-isma/app"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	a, err := app.Build()
	if err != nil {
		return err
	}
	a.Bot.Start()
	return nil
}
