package main

import (
	"context"
	"english-speaking-club-bot/internal/config"
	"english-speaking-club-bot/internal/server"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if os.Getenv("RAILWAY_ENVIRONMENT") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("error loading environment variables")
			return
		}
	}
	fmt.Println("Start")

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal("error loading config")
		return
	}

	srv := server.NewServer(conf)

	if err := srv.Start(); err != nil {
		log.Fatal("error starting server" + err.Error())
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()
	err = srv.Stop()

	fmt.Println("Shutting down...")
	if err != nil {
		fmt.Println(err.Error())
	}
}
