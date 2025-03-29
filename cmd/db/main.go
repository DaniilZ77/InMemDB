package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/DaniilZ77/InMemDB/internal/app"
	"github.com/DaniilZ77/InMemDB/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config := config.MustConfig()
	app, err := app.NewApp(ctx, config)
	if err != nil {
		panic(err)
	}

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	app.Shutdown(ctx)
}
