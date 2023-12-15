package models

//Структура товара.
type Item struct {
	ChrtId      int    `json:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int    `json:"price" validate:"required,gt=0"`
	Rid         string `json:"rid" validate:"required,min=21,max=21"`
	Name        string `json:"name" validate:"required"`
	Sale        int    `json:"sale" validate:"gte=0,lte=100"`
	Size        string `json:"size" validate:"required"`
	TotalPrice  int    `json:"total_price" validate:"required, gte=0"`
	NmId        int    `json:"nm_id" validate:"required"`
	Brand       string `json:"brand" validate:"required"`
	Status      int    `json:"status" validate:"required"`
}
