package main

import (
	"context"

	"github.com/theskyinflames/quiz/csvreader/config"
	"github.com/theskyinflames/quiz/csvreader/reader"
	"github.com/theskyinflames/quiz/csvreader/repository"
	"github.com/theskyinflames/quiz/csvreader/service"
)

func main() {
	ctx := context.Background()

	cfg := &config.Config{}
	cfg.Load()

	repository := repository.NewRepository(cfg)
	err := repository.Connect()
	if err != nil {
		panic(err)
	}

	csvService := service.NewService(ctx, cfg, repository)
	err = csvService.Start(reader.DefaultCSVReaderFunc)
	if err != nil {
		panic(err)
	}
}
