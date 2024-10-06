package main

import (
	goGO "github.com/goGo-service/back"
	"github.com/goGo-service/back/pkg/handler"
	"github.com/goGo-service/back/pkg/repository"
	"github.com/goGo-service/back/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title goGO
// @version1.0
// @description goGo-service

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}
	// TODO: выпилить godotenv
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		Username: viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASS"),
		DBName:   viper.GetString("DB_NAME"),
		SSLMode:  viper.GetString("DB_SSLMODE"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(goGO.Server)

	if err := srv.Run(viper.GetString("HOST_PORT"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}

	logrus.Print("goGO started")
}

func initConfig() error {
	viper.AutomaticEnv()
	return nil
}
