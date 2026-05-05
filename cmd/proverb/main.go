package main

import (
	"log"

	"github.com/BogdanBratsky/proverb/internal/app"
	"github.com/BogdanBratsky/proverb/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	srv, err := app.NewServer(&cfg.AppConfig)
	if err != nil {
		panic(err)
	}

	log.Println("server is running")
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
