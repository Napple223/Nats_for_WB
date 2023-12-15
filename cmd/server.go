package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"wb-lvl0/internal/api"
	"wb-lvl0/internal/app/consumer"
	"wb-lvl0/internal/app/producer"
	"wb-lvl0/internal/models"
	"wb-lvl0/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	//Загружаем переменные окружения и парсим их в соответсвующие структуры.
	err := godotenv.Load("./config.env")
	if err != nil {
		log.Fatal("файл конфигурации не найден")
	}
	d := models.ConfigDB{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
		UserName: os.Getenv("USER_NAME"),
		Password: os.Getenv("PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
		SslMode:  os.Getenv("SSL_MODE"),
	}
	nPub := models.ConfigNatsPub{
		NatsURL:        os.Getenv("NATS_URL_PUB"),
		NatsCluster:    os.Getenv("CLUSTER_ID"),
		NatsClient:     os.Getenv("CLIENT_PUB"),
		NatsSubjectPub: os.Getenv("NATS_SUBJECT"),
	}
	nSub := models.ConfigNatsSub{
		NatsURL:        os.Getenv("NATS_URL_SUB"),
		NatsCluster:    os.Getenv("CLUSTER_ID"),
		NatsClient:     os.Getenv("CLIENT_SUB"),
		NatsSubjectSub: os.Getenv("NATS_SUBJECT"),
	}
	done := make(chan os.Signal, 1)
	var wg sync.WaitGroup
	//Инициализиуем новый экземпляр хранилища (кэш + БД).
	s, err := storage.New(d)
	if err != nil {
		log.Fatal(err)
	}

	//Инициализиуем продюссера данных nats streaming.
	go producer.Run(nPub, &wg, done)
	//Инициализируем подписчика nats streaming.
	go consumer.Run(nSub, &wg, s)
	//Инициализируем api.
	api := api.New(s)
	//Запускаем сервер.
	err = http.ListenAndServe(os.Getenv("APP_PORT"), api.Router())

	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)
	close(done)
	wg.Wait()
}
