package main

import (
	"context"
	"fmt"
	"net/http"
	"nutritiontracker/handler"
	"nutritiontracker/mongo"
	"os"
	"os/signal"
)

type Main struct {
	DB         *mongo.DB
	HTTPServer *handler.Server
}

func New() *Main {
	db := mongo.NewDB("mongodb://localhost:27017/")
	return &Main{
		DB:         db,
		HTTPServer: handler.NewServer(),
	}
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

	m := New()

	if err := m.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	<-ctx.Done()
}
