package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"workmate/internal/handler"
	"workmate/internal/repository"
	"workmate/internal/service"

	"github.com/sirupsen/logrus"
)

const (
	_shutdownPeriod      = 15 * time.Second
	_shutdownHardPeriod  = 3 * time.Second
	_readinessDrainDelay = 5 * time.Second
)

func main() {
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	repo := repository.NewRepository()
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	app := handlers.InitRoutes()

	var once sync.Once
	once.Do(func() {
		logrus.Info("Init persistent storage")
		err := repo.InitPersistentStorage()
		if err != nil {
			logrus.Fatal(err)
		}
	})

	go func() {
		logrus.Info("starting server")
		if err := app.Listen(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("error starting server: %v", err)
		}
	}()

	// крон для обновления данных в персистентном хранилище
	go func() {
		for {
			select {
			case <-time.After(time.Second * 5):
				err := repo.URLStorageRepository.WriteDataToFileAndLocalStorage()
				if err != nil {
					logrus.Errorf("error storing data: %v", err)
				}
			case <-rootCtx.Done():
				err := repo.URLStorageRepository.WriteDataToFileAndLocalStorage()
				if err != nil {
					logrus.Errorf("error storing data: %v", err)
				}
				return
			}
		}
	}()

	<-rootCtx.Done()
	stop()

	handler.IsShutDown.Store(true)
	logrus.Info("Received shutdown signal, shutting down.")

	time.Sleep(_readinessDrainDelay)
	logrus.Info("Readiness check propagated, now waiting for ongoing requests to finish.")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), _shutdownPeriod)
	defer cancel()

	err := app.ShutdownWithContext(shutdownCtx)
	if err != nil {
		logrus.Info("Failed to wait for ongoing requests to finish, waiting for forced cancellation.")
		time.Sleep(_shutdownHardPeriod)
	}
}
