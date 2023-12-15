package cache

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"wb-lvl0/internal/models"
	"wb-lvl0/internal/storage/postgres"
)

// Структура кэша.
type Cache struct {
	db   *postgres.DB
	m    sync.RWMutex
	Data map[string]models.Order
}

// Функция-конструктор для создания нового экземпляра кэша, принимает в себя конфиг
// для подключения к БД.
func New(cfg models.ConfigDB) (*Cache, error) {
	dsn := fmt.Sprintf("%s://%s:%s@%s/%s",
		cfg.Host,
		cfg.UserName,
		cfg.Password,
		cfg.Port,
		cfg.DbName)
	db, err := postgres.New(dsn)
	if err != nil {
		return nil, err
	}
	d := make(map[string]models.Order)
	orders, err := db.GetAll() //тут могут быть проблемы в случае большого объема данных в БД.
	if err != nil {
		log.Println(err)
		return nil, err
	} else {
		for _, v := range orders {
			d[v.OrderUid] = v
		}
	}
	c := Cache{
		db:   db,
		Data: d,
	}

	return &c, nil
}

// Функция для добавления заказа в кэш и в БД.
func (c *Cache) Add(o models.Order) error {
	_, ok := c.Data[o.OrderUid]
	if ok {
		return errors.New("найден дубликат")
	}
	err := c.db.Add(o)
	if err != nil {
		return err
	}
	c.Data[o.OrderUid] = o
	return nil
}

// Функция для получения заказа по uid заказа.
func (c *Cache) Get(uid string) (models.Order, error) {
	v, ok := c.Data[uid]
	if ok {
		return v, nil
	}
	err := fmt.Sprintf("заказ с uid %s не найден\n", uid)
	return v, errors.New(err)
}
