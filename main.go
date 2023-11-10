// Package dutch-bot generates Dutch grammar rules using ChatGPT and
// publishes them to Telegram.
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	configFile := flag.String("config", "./config.yaml", "path to config file")
	runOnce := flag.String("once", "", "generate one text and exit")
	debug := flag.Bool("debug", false, "run in debug mode")
	flag.Parse()

	conf, err := ReadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	conf.Debug = conf.Debug || *debug

	log.Print("Starting...")

	// Create and run bot
	bot, err := NewBot(conf)
	if err != nil {
		log.Fatalf("Failed to init bot: %v", err)
	}
	if *runOnce != "" {
		once(ctx, bot, *runOnce)
	} else {
		loop(ctx, bot)
	}
	log.Print("Shutdown")
}

func once(ctx context.Context, bot *Bot, promptName string) {
	if err := bot.RunOnce(ctx, promptName); err != nil {
		log.Printf("Bot failed with error: %v", err)
	}
}

func loop(ctx context.Context, bot *Bot) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		bot.Run(ctx)
		wg.Done()
	}()
	<-ctx.Done() // wait for SIGTERM/SIGINT
	wg.Wait()
}
