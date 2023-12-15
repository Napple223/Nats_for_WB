package producer

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
	"wb-lvl0/internal/models"
	"wb-lvl0/internal/nats"
)

const (
	period int = 15
)

func Run(nPub models.ConfigNatsPub, wg *sync.WaitGroup, done <-chan os.Signal) {
	pub, err := nats.New(nPub.NatsCluster, nPub.NatsClient, nPub.NatsURL)
	if err != nil {
		log.Fatalf("ошибка подключения продюсера %v", err)
	}
	data := make(chan []byte)
	wg.Add(2)
	go dataSourse(done, data, wg)
	go pubHandler(done, data, pub, nPub.NatsSubjectPub, wg)
	wg.Wait()
	close(data)
}

func dataSourse(done <-chan os.Signal, data chan<- []byte, wg *sync.WaitGroup) {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	order := models.Order{
		OrderUid:    "b563feb7b2b84b6",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: models.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestId:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []models.Item{{
			ChrtId:      9934930,
			TrackNumber: "WBILMTESTTRACK",
			Price:       453,
			Rid:         "ab4219087a764ae0btest",
			Name:        "Mascaras",
			Sale:        30,
			Size:        "0",
			TotalPrice:  317,
			NmId:        2389212,
			Brand:       "Vivienne Sabo",
			Status:      202,
		}},
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmID:              99,
		DateCreated:       "2021-11-26T06:22:19Z",
		OffShard:          "1",
	}
	loop := true
	for loop {
		select {
		case <-done:
			loop = false
		default:
			i := rand.Intn(9000) + 1000
			or := "b563feb7b2b84b6"
			order.OrderUid = or + strconv.Itoa(i)
			dataByte, err := json.Marshal(order)
			if err != nil {
				log.Println("ошибка маршализации json")
			}
			data <- dataByte
			time.Sleep(time.Duration(period) * time.Second)
		}
	}
	wg.Done()
}

func pubHandler(done <-chan os.Signal, data <-chan []byte, n *nats.Nats, subject string, wg *sync.WaitGroup) {
	loop := true
	for loop {
		select {
		case <-done:
			log.Println("остановка отправки сообщений")
			loop = false
		case d := <-data:
			err := n.Conn.Publish(subject, d)
			if err != nil {
				log.Printf("ошибка публикации данных в натс:%v", err)
			}
		default:
		}
	}
	wg.Done()
}
