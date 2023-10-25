package main

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	wb_l0 "wb-l0"
	"wb-l0/pkg/handlers"
	"wb-l0/pkg/repository"
	"wb-l0/pkg/service"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("Ошибка при инициализации конфигурации: %s", err.Error())
	}
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBname:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: viper.GetString("db.password"),
	})
	if err != nil {
		logrus.Fatalf("Ошибка при инициализации базы данных: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	err = services.PullAllOrders()
	if err != nil {
		logrus.Fatalf("Ошибка при восстановлении кэша: %s", err.Error())
	}
	handlers := handlers.NewHandler(services)
	srv := new(wb_l0.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Ошибка при запуске сервера: %s", err.Error())
		}
	}()

	nc, err := nats.Connect(viper.GetString("nats.url"))
	if err != nil {
		logrus.Fatalf("Ошибка подключения к nats-streaming: %s", err.Error())
	}
	defer nc.Close()
	channel := viper.GetString("nats.channel")
	sub, err := nc.Subscribe(channel, func(msg *nats.Msg) {
		logrus.Printf("Получено сообщение из канала: %s\n", string(msg.Data))
		handlers.AddOrder(msg.Data)
	})
	if err != nil {
		logrus.Errorf("Ошибка при подписки на канал %s: %s", channel, err.Error())
	}
	defer sub.Unsubscribe()

	logrus.Println("Запуск сервиса")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("Завершение сервиса")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Ошибка при завершении работы сервера: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("Ошибка при закрытии соединения с БД: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
