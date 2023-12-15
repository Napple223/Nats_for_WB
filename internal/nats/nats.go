package nats

import (
	"log"

	"github.com/nats-io/stan.go"
)

// Объект натс
type Nats struct {
	Conn stan.Conn
}

// Функция конструктор, возвращающая подключение к натс.
func New(natsCluster, natsClient, natsURL string) (*Nats, error) {
	conn, err := stan.Connect(natsCluster, natsClient, stan.NatsURL(natsURL))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	n := Nats{
		Conn: conn,
	}
	return &n, nil
}
