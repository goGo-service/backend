package main

import (
	"fmt"
	"log"
	"os"

	goGO "github.com/goGo-service/back"
	"github.com/goGo-service/back/pkg/handler"
	"github.com/goGo-service/back/pkg/repository"
	"github.com/goGo-service/back/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	// "github.com/spf13/viper"
)

func main() {
	// if err := initConfig(); err != nil {
	// 	log.Fatalf("error initializing configs: %s", err.Error())
	// }
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}
	fmt.Println(os.Getenv("DB_USERNAME"))

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(goGO.Server)
	if err := srv.Run(os.Getenv("PORT"), handlers.InitRoutes()); err != nil {
		log.Panic()
	}
}

// func initConfig() error {
// 	viper.SetConfigName(".env")
// 	return viper.ReadInConfig()
// }
