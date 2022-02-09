package main

import (
	"chainrunner/flags"
	"chainrunner/internal/bot"
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	flag.Parse()

	// cpu profiling
	if *flags.CPU_PROFILE != "" {
		f, err := os.Create(*flags.CPU_PROFILE)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
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

	if *flags.MEM_PROFILE != "" {
		f, err := os.Create(*flags.MEM_PROFILE)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
