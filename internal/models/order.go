package models

import (
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Order struct {
	OrderUID          string    `json:"order_uid" db:"order_uid"`
	TrackNumber       string    `json:"track_number" db:"track_number"`
	Entry             string    `json:"entry" db:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale" db:"locale"`
	InternalSignature string    `json:"internal_signature" db:"internal_signature"`
	CustomerId        string    `json:"customer_id" db:"customer_id"`
	DeliveryService   string    `json:"delivery_service" db:"delivery_service"`
	Shardkey          string    `json:"shardkey" db:"shardkey"`
	SmId              int       `json:"sm_id" db:"sm_id"`
	DateCreated       time.Time `json:"date_created" db:"date_created"`
	OofShard          string    `json:"oof_shard" db:"oof_shard"`
	InDB              bool      `json:"-"`
}

func (o *Order) Validate() error {
	return validation.ValidateStruct(
		o,
		validation.Field(&o.OrderUID, validation.Required),
		validation.Field(&o.TrackNumber, validation.Required),
		validation.Field(&o.Entry, validation.Required),
		validation.Field(&o.Locale, validation.Required),
		validation.Field(&o.InternalSignature, validation.Required),
		validation.Field(&o.CustomerId, validation.Required),
		validation.Field(&o.DeliveryService, validation.Required),
		validation.Field(&o.Shardkey, validation.Required),
		validation.Field(&o.SmId, validation.Required, is.Int),
		validation.Field(&o.DateCreated, validation.Required),
		validation.Field(&o.OofShard, validation.Required),
	)
}

func (o *Order) IsEmpty() bool {
	return o.OrderUID == ""
}
