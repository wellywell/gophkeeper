package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"net/http"
	_ "net/http/pprof"

	"github.com/wellywell/gophkeeper/internal/config"
	"github.com/wellywell/gophkeeper/internal/db"
	"github.com/wellywell/gophkeeper/internal/handlers"
	"github.com/wellywell/gophkeeper/internal/logging"
	"github.com/wellywell/gophkeeper/internal/router"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s", buildVersion, buildDate, buildCommit)

	logger, err := logging.NewLogger()
	if err != nil {
		panic(err)
	}

	conf, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	database, err := db.NewDatabase(conf.DatabaseDSN)

	if err != nil {
		panic(err)
	}
	defer func() {
		err = database.Close()
		if err != nil {
			panic(err)
		}
	}()

	hndl := handlers.NewHandlerSet(conf.Secret, database)

	s := router.NewServer(*conf, *hndl, logger)

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	go func() {
		<-sig
		// Trigger graceful shutdown
		err := s.Shutdown(serverCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	err = s.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
