// Package main is entry point for port domain service
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	ps "github.com/dehimb/shrimp/port-domain-service/internal/port-domain-service"
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
	ps.StartPortDomainService(ctx, logger)
}
