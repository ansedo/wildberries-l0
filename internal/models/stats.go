package models

type Stats struct {
	OrderCount    int      `json:"order_count"`
	ItemCount     int      `json:"item_count"`
	LastOrderUIDs []string `json:"last_order_uids"`
}
