package models

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Payment struct {
	Transaction  string `json:"transaction" db:"transaction"`
	RequestId    string `json:"request_id" db:"request_id"`
	Currency     string `json:"currency" db:"currency"`
	Provider     string `json:"provider" db:"provider"`
	Amount       int    `json:"amount" db:"amount"`
	PaymentDt    int    `json:"payment_dt" db:"payment_dt"`
	Bank         string `json:"bank" db:"bank"`
	DeliveryCost int    `json:"delivery_cost" db:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total" db:"goods_total"`
	CustomFee    int    `json:"custom_fee" db:"custom_fee"`
}

func (p *Payment) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(&p.Transaction, validation.Required),
		validation.Field(&p.RequestId, validation.Required),
		validation.Field(&p.Currency, validation.Required),
		validation.Field(&p.Provider, validation.Required),
		validation.Field(&p.Amount, validation.Required, is.Int),
		validation.Field(&p.PaymentDt, validation.Required, is.Int),
		validation.Field(&p.Bank, validation.Required),
		validation.Field(&p.DeliveryCost, validation.Required, is.Int),
		validation.Field(&p.GoodsTotal, validation.Required, is.Int),
		validation.Field(&p.CustomFee, validation.Required, is.Int),
	)
}
