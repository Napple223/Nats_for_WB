//Пакет, реализующий хранение заказов в БД и их подгрузку в кэщ.

package postgres

import (
	"testing"
	"wb-lvl0/internal/models"
)

func TestDB(t *testing.T) {
	dsn := "postgres://andrey:andrey@localhost:5432/WB"
	db, err := New(dsn)
	if err != nil {
		t.Fatal(err)
	}

	i := models.Item{
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
	}
	want := models.Order{
		OrderUid:    "b563feb7b2b84b61111",
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
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmID:              99,
		DateCreated:       "2021-11-26T06:22:19Z",
		OffShard:          "1",
		Items:             make([]models.Item, 1, 1),
	}
	want.Items[0] = i
	err = db.Add(want)
	if err != nil {
		t.Error(err)
	}
	got, err := db.GetAll()
	if err != nil {
		t.Error(err)
	}
	if got[0].OrderUid != want.OrderUid {
		t.Errorf("получили %v, а ожидали %v", got[0].OrderUid, want.OrderUid)
	}
}
