package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/DaniilZ77/InMemDB/internal/app"
	"github.com/DaniilZ77/InMemDB/internal/config"
	_ "go.uber.org/automaxprocs"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config := config.MustConfig()
	if err := app.RunApp(ctx, config); err != nil {
		panic(err)
	}
}
