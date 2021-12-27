package main

import (
	"chainrunner/flags"
	"chainrunner/internal/bot"
	"flag"
	"log"
	"os"

	"github.com/pkg/profile"
)

func main() {
	flag.Parse()

	if *flags.ENABLE_PPROF {
		defer profile.Start().Stop()
	}

	bot, err := bot.NewBot(*flags.DB_PATH, *flags.DEFAULT_CLIENT)

	if err != nil {
		log.Print("something wrong with making a new bot ", err)
		os.Exit(-1)
	}

	if err := bot.Run(); err != nil {
		log.Println("something wrong with running the bot", err)
	}

	if err := bot.CloseResources(); err != nil {
		log.Print("something wrong with closing resources", err)
	}
}