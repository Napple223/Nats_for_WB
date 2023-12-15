package consumer

import (
	"encoding/json"
	"log"
	"sync"
	"wb-lvl0/internal/models"
	"wb-lvl0/internal/nats"
	"wb-lvl0/internal/storage"
	"wb-lvl0/internal/validator"

	"github.com/nats-io/stan.go"
)

func Run(nSub models.ConfigNatsSub, wg *sync.WaitGroup, s *storage.Storage) {
	v := validator.New()
	sub, err := nats.New(nSub.NatsCluster, nSub.NatsClient, nSub.NatsURL)
	if err != nil {
		log.Fatalf("не удалось подписаться на натс %v:", err)
	}
	wg.Add(1)
	go subHandler(sub, s, v, nSub.NatsSubjectSub, wg)
	wg.Wait()
}

func subHandler(n *nats.Nats, s *storage.Storage, v *validator.Validator, subject string, wg *sync.WaitGroup) {
	sub, err := n.Conn.Subscribe(subject, func(msg *stan.Msg) {
		var order models.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Println("не удалось расшифровать сообщение из натс")
		}

		if v.ValidateStruct(order) {
			s.Add(order)
		}
	})

	for sub.IsValid() {

	}
	err = sub.Unsubscribe()
	if err != nil {
		log.Println("не удалось отписаться от натс")
	}
	wg.Done()
}
