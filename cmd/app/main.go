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
			log.Fatal("error loading environment variables" + err.Error())
			return
		}
	}
	fmt.Println("Start")

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal("error loading config" + err.Error())
		return
	}

	srv := server.NewServer(conf)

	if err := srv.Start(); err != nil {
		log.Fatal("error starting server" + err.Error())
	}

	//voc, err := services.LoadVocabulary("common_words.txt")
	//if err != nil {
	//	log.Fatal("error loading vocabulary" + err.Error())
	//}
	//
	//w1 := services.GetRandomWordFromVocabulary(voc)
	//
	//word, err := services.FetchWordWithTranslation(os.Getenv("YANDEX_DICT_API_KEY"), w1)
	//if err != nil {
	//	fmt.Println("error getting word")
	//} else {
	//	fmt.Println(word)
	//}

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
