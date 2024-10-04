package main

import (
	"fmt"
	goGO "github.com/goGo-service/back"
	"github.com/goGo-service/back/pkg/handler"
	"github.com/goGo-service/back/pkg/repository"
	"github.com/goGo-service/back/pkg/repository/cache"
	"github.com/goGo-service/back/pkg/service"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
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
	//ctx := context.Background()

	redisClient, err := cache.NewRedisDB()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	defer redisClient.Close()

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

const (
	vkIdUsersTable  = "vk_id_users"
	UsersTable      = "users"
	roomsTable      = "rooms"
	rolesTable      = "roles"
	roomsUsersTable = "rooms_users"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func initConfig() error {
	viper.AutomaticEnv()
	return nil
}
