package main

import (
	"log"

	goGO "github.com/goGo-service/back"
	"github.com/goGo-service/back/pkg/handler"
	"github.com/goGo-service/back/pkg/repository"
	"github.com/goGo-service/back/pkg/service"
)

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(goGO.Server)
	if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
		log.Panic()
	}
}
