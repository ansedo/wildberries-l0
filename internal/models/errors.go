package models

import "errors"

var (
	ErrOrderUIDRequired = errors.New("order uid field is required")
	ErrOrderUIDNotExist = errors.New("this order uid does not exist")
)
