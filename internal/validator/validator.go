package validator

import (
	"log"
	"wb-lvl0/internal/models"

	"github.com/go-playground/validator/v10"
)

// Структура валидатора.
type Validator struct {
	validate *validator.Validate
}

// Функция конструктор для создания нового экземпляра валидатора.
func New() *Validator {
	v := validator.New(validator.WithRequiredStructEnabled())
	vld := Validator{
		validate: v,
	}
	return &vld
}

// Функция проверки валидности структуры заказа.
func (v *Validator) ValidateStruct(o models.Order) bool {
	err := v.validate.Struct(o)
	if err != nil {
		log.Printf("ошибка валидатора, %v", err)
		return false
	}
	return true
}
