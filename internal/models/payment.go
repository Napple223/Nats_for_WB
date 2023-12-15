package models

//Структура оплаты.
type Payment struct {
	Transaction  string `json:"transaction" validate:"required,min=19,max=19"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency" validate:"required"`
	Provider     string `json:"provider" validate:"required"`
	Amount       int    `json:"amount" validate:"required,gt=0"`
	PaymentDt    int    `json:"payment_dt" validate:"required"`
	Bank         string `json:"bank" validate:"required"`
	DeliveryCost int    `json:"delivery_cost" validate:"required,gte=0"`
	GoodsTotal   int    `json:"goods_total" validate:"required,gt=0"`
	CustomFee    int    `json:"custom_fee"`
}
