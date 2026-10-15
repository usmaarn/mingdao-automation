package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/usmaarn/mingdao/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	app := app.New(ctx)
	app.OpenFormModal()

	server := &http.Server{
		Handler:      app,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		Addr:         ":9922",
	}

	go func() {
		log.Println("Server running on", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Println("Server error: ", err)
			}
		}
	}()

	<-ctx.Done()
	log.Println("shutting down ....")
	err := server.Shutdown(ctx)
	if err != nil {
		err := server.Close()
		if err != nil {
			os.Exit(1)
		}
	}

	log.Println("Server shutdoen completed.")
}
