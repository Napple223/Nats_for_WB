package storage

import (
	"wb-lvl0/internal/models"
	"wb-lvl0/internal/storage/cache"
)

// структура хранилища заказов.
type Storage struct {
	Cache
}

type Cache interface {
	Add(order models.Order) error
	Get(uid string) (models.Order, error)
}

// Функция конструктор для создания нового экземпляра хранилища.
func New(cfg models.ConfigDB) (*Storage, error) {
	c, err := cache.New(cfg)
	if err != nil {
		return nil, err
	}
	s := Storage{
		Cache: c,
	}
	return &s, nil
}
