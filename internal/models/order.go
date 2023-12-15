package models

//Структура заказа.
type Order struct {
	OrderUid          string   `json:"order_uid" validate:"required,min=19,max=19"`
	TrackNumber       string   `json:"track_number" validate:"required,min=14,max=14"`
	Entry             string   `json:"entry" validate:"required"`
	Delivery          Delivery `json:"delivery" validate:"required"`
	Payment           Payment  `json:"payment" validate:"required"`
	Items             []Item   `json:"items" validate:"required"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerId        string   `json:"customer_id" validate:"required"`
	DeliveryService   string   `json:"delivery_service" validate:"required"`
	ShardKey          string   `json:"shardkey"`
	SmID              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created" validate:"required"`
	OffShard          string   `json:"off_shard"`
}
