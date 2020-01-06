package clientapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	port = ":8080"
)

func StartServer(ctx context.Context, logger *logrus.Logger) {
	handler := &handler{
		router: mux.NewRouter(),
		logger: logger,
	}
	s := &http.Server{
		Addr:         port,
		Handler:      handler,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	err := handler.initRouter(ctx)
	if err != nil {
		logger.Fatal("Can't initialize router: ", err)
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("Error starting server", err)
		}
	}()

	logger.Print("Started server on port ", port)
	waitForShutdown(ctx, s, logger)
}

func waitForShutdown(ctx context.Context, s *http.Server, logger *logrus.Logger) {
	<-ctx.Done()
	logger.Info("Trying graceful shutdown server")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctxShutDown); err != nil {
		logger.Errorf("Server shutdown failed: %s", err)
		return
	}
	logger.Info("Server stopped")
}
