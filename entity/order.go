package entity

import "time"

type Order struct {
	Order_id      int       `json:"order_id"`
	Customer_name string    `json:"customer_name"`
	Ordered_at    time.Time `json:"ordered_at"`
	Item          []Item    `json:"item"`
}

type Item struct {
	Item_id     int    `json:"item_id"`
	Item_code   string `json:"item_code"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Order_id    int    `json:"order_id"`
}
