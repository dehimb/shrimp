// Package main is entry point for client api service
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dehimb/shrimp/client-api/internal/clientapi"
	"github.com/dehimb/shrimp/client-api/internal/portsloader"
	"github.com/sirupsen/logrus"
)

func main() {
	// Create logger instance
	logger := logrus.New()
	// Set logger leve according to environment
	switch os.Getenv("APP_ENV") {
	case "prod":
		logger.SetLevel(logrus.InfoLevel)
	default:
		logger.SetLevel(logrus.DebugLevel)
	}

	// Catch syscal signals to graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := <-c
		logger.Infof("Syscall signal: %s", sig)
		cancel()
	}()

	loadPortsData(ctx, logger)

	// Start client api service
	clientapi.StartServer(ctx, logger)
}

func loadPortsData(ctx context.Context, logger *logrus.Logger) {
	// Trying to read file with ports
	file, err := os.Open("./data/ports.json")
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infoln("Updating ports database")
	results, err := portsloader.Load(ctx, file)
	if err != nil {
		logger.Warnln("Error when updating ports data: ", err)
	}
	if results != nil {
		logger.Infof("Ports data updated: %d records", results.Count)
	}
}
