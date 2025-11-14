package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"sync/atomic"
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

var isShutDown atomic.Bool

func main() {
	// Отлавливаем системные сигналы о запланированном завершении работы текущего процесса ОС
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Объявляем нижний слой "repository"
	repo := repository.NewRepository()

	services := service.NewService(repo)

	// Объявляем верхний слой "handler"
	handlers := handler.NewHandler(services)

	// Объявляем экземпляр нашего сервера
	app := handlers.InitRoutes(&isShutDown)

	// Сервер запускаем в отдельной горутине, поскольку app.Listen не будет работать с app.ShutdownWithContext в рамках одной горутины
	go func() {
		logrus.Info("starting server")
		if err := app.Listen(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("error starting server: %v", err)
		}
	}()

	// Ждём ответ от канала контекста системного завершения процесса
	<-rootCtx.Done()

	// Прокидываем уведомление об аварийном уничтожении процесса
	isShutDown.Store(true)
	logrus.Println("Received shutdown signal, shutting down.")

	// Give time for readiness check to propagate
	time.Sleep(_readinessDrainDelay)
	logrus.Println("Readiness check propagated, now waiting for ongoing requests to finish.")

	// Создаём контекст с таймером, чтобы дать запущенным горутинам(хэндлерам) время для завершения своей работы
	shutdownCtx, cancel := context.WithTimeout(context.Background(), _shutdownPeriod)
	defer cancel()

	err := app.ShutdownWithContext(shutdownCtx)
	if err != nil {
		logrus.Println("Failed to wait for ongoing requests to finish, waiting for forced cancellation.")
		time.Sleep(_shutdownHardPeriod)
	}
}
