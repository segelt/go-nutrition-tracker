package main

import (
	"context"
	"fmt"
	"net/http"
	"nutritiontracker/config"
	"nutritiontracker/handler"
	"nutritiontracker/mongo"
	"os"
	"os/signal"

	"github.com/go-playground/validator"
)

type Main struct {
	DB         *mongo.DB
	HTTPServer *handler.Server
}

func New() (*Main, error) {
	conf, err := config.NewMealConfig()
	if err != nil {
		return nil, fmt.Errorf("Config.New :%s", err)
	}

	validate := validator.New()
	if err := validate.Struct(conf); err != nil {
		return nil, fmt.Errorf("Configuration.Validate: Missing required attributes %v", err)
	}

	return &Main{
		DB:         mongo.NewDB(conf.BaseDBConfig),
		HTTPServer: handler.NewServer(conf.BaseServerConfig),
	}, nil
}

func (m *Main) Run(ctx context.Context) (err error) {
	if err := m.DB.Open(); err != nil {
		return fmt.Errorf("cannot open db: %w", err)
	}

	mealService := mongo.NewMealService(m.DB)
	m.HTTPServer.MealService = mealService

	if err := m.HTTPServer.Start(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Server startup failure: %s", err.Error())
	}

	return nil
}

// Close gracefully stops the program.
func (m *Main) Close() error {
	if m.HTTPServer != nil {
		if err := m.HTTPServer.Close(); err != nil {
			return err
		}
	}
	if m.DB != nil {
		if err := m.DB.Close(); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()

	m, err := New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Main.New() %v\n", err)
		os.Exit(1)
	}

	if err := m.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	<-ctx.Done()
}
