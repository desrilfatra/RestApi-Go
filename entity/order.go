package entity

import "time"

type Order struct {
	Order_id      int       `json:"order_id"`
	Customer_name string    `json:"customerName"`
	Ordered_at    time.Time `json:"orderedAt"`
	Item          []Item    `json:"items"`
}

type Item struct {
	Item_id     int    `json:"lineItemId"`
	Item_code   string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Order_id    int    `json:"order_id"`
}
