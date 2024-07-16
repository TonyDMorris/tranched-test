package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/tonydmorris/tranched/internal/app"
	"github.com/tonydmorris/tranched/internal/repository/order"
	"github.com/tonydmorris/tranched/internal/repository/user"

	"github.com/caarlos0/env/v10"
	_ "github.com/lib/pq"
	"github.com/tonydmorris/tranched/pkg/logging"
)

type Config struct {
	PostgresHost     string `env:"POSTGRES_HOST,required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresUser     string `env:"POSTGRES_USER,required"`
	PostgresDB       string `env:"POSTGRES_DB,required"`
	PostgresPort     int    `env:"POSTGRES_PORT" envDefault:"5432"`
}

func main() {
	var cfg Config
	// Load configuration
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	logger, err := logging.NewProductionWithSugar()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	defer db.Close()

	userRepo := user.NewRepository(db)
	orderRepo := order.NewRepository()

	api := app.New(userRepo, orderRepo, app.WithLogger(logger))

	err = api.Run()
	if err != nil {
		log.Fatalf("failed to run api: %v", err)
	}

}
