package models

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Delivery struct {
	OrderUID string `json:"-" db:"order_uid"`
	Name     string `json:"name" db:"name"`
	Phone    string `json:"phone" db:"phone"`
	Zip      string `json:"zip" db:"zip"`
	City     string `json:"city" db:"city"`
	Address  string `json:"address" db:"address"`
	Region   string `json:"region" db:"region"`
	Email    string `json:"email" db:"email"`
}

func (d *Delivery) Validate() error {
	return validation.ValidateStruct(
		d,
		validation.Field(&d.Name, validation.Required),
		validation.Field(&d.Phone, validation.Required),
		validation.Field(&d.Zip, validation.Required),
		validation.Field(&d.City, validation.Required),
		validation.Field(&d.Address, validation.Required),
		validation.Field(&d.Region, validation.Required),
		validation.Field(&d.Email, validation.Required, is.EmailFormat),
	)
}
