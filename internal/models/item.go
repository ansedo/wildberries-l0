package models

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Item struct {
	ChrtId      int    `json:"chrt_id" db:"chrt_id"`
	TrackNumber string `json:"track_number" db:"track_number"`
	Price       int    `json:"price" db:"price"`
	Rid         string `json:"rid" db:"rid"`
	Name        string `json:"name" db:"name"`
	Sale        int    `json:"sale" db:"sale"`
	Size        string `json:"size" db:"size"`
	TotalPrice  int    `json:"total_price" db:"total_price"`
	NmId        int    `json:"nm_id" db:"nm_id"`
	Brand       string `json:"brand" db:"brand"`
	Status      int    `json:"status" db:"status"`
}

func (i *Item) Validate() error {
	return validation.ValidateStruct(
		i,
		validation.Field(&i.ChrtId, validation.Required, is.Int),
		validation.Field(&i.TrackNumber, validation.Required),
		validation.Field(&i.Price, validation.Required, is.Int),
		validation.Field(&i.Rid, validation.Required),
		validation.Field(&i.Name, validation.Required),
		validation.Field(&i.Sale, validation.Required, is.Int),
		validation.Field(&i.Size, validation.Required),
		validation.Field(&i.TotalPrice, validation.Required, is.Int),
		validation.Field(&i.NmId, validation.Required, is.Int),
		validation.Field(&i.Brand, validation.Required),
		validation.Field(&i.Status, validation.Required, is.Int),
	)
}
