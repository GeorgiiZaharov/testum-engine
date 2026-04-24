package main

import (
	"log"

	"testum-engine/app/internal/app"
	"testum-engine/app/internal/config"
)

func main() {
	cfg := config.Load()

	application, err := app.Build(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
