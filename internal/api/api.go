package api

import (
	"encoding/json"
	"net/http"
	"wb-lvl0/internal/storage"

	"github.com/gorilla/mux"
)

// Объект программного интерфейса сервера
type API struct {
	s      storage.Storage
	router *mux.Router
}

// Регистрация api в маршрутизаторе запросов.
func (api *API) endpoints() {
	api.router.Use(api.Middleware)
	api.router.HandleFunc("/orders/{uid}", api.ordersHandler).Methods(http.MethodGet, http.MethodOptions)
}

// Функция-конструктор объекта API.
func New(s *storage.Storage) *API {
	api := API{
		s: *s,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// Получение маршрутизатора запросов.
func (api *API) Router() *mux.Router {
	return api.router
}

// Функция для получения заказа по uid.
func (api *API) ordersHandler(res http.ResponseWriter, req *http.Request) {
	uid := mux.Vars(req)["uid"]
	order, err := api.s.Cache.Get(uid)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(res).Encode(order)
}
