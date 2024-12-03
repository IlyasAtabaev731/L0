package models

import (
	"time"
)

// Order represents the orders table
type Order struct {
	OrderUid          string    `json:"order_uid" storage:"order_uid"`
	TrackNumber       string    `json:"track_number" storage:"track_number"`
	Entry             string    `json:"entry" storage:"entry"`
	Locale            string    `json:"locale" storage:"locale"`
	InternalSignature string    `json:"internal_signature" storage:"internal_signature"`
	CustomerID        string    `json:"customer_id" storage:"customer_id"`
	DeliveryService   string    `json:"delivery_service" storage:"delivery_service"`
	Shardkey          string    `json:"shardkey" storage:"shardkey"`
	SMID              int       `json:"sm_id" storage:"sm_id"`
	DateCreated       time.Time `json:"date_created" storage:"date_created"`
	OofShard          string    `json:"oof_shard" storage:"oof_shard"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
}

// Delivery represents the deliveries table
type Delivery struct {
	OrderUid string `json:"order_uid" storage:"order_uid"`
	Name     string `json:"name" storage:"name"`
	Phone    string `json:"phone" storage:"phone"`
	Zip      string `json:"zip" storage:"zip"`
	City     string `json:"city" storage:"city"`
	Address  string `json:"address" storage:"address"`
	Region   string `json:"region" storage:"region"`
	Email    string `json:"email" storage:"email"`
}

// Payment represents the payments table
type Payment struct {
	OrderUid      string    `json:"order_uid" storage:"order_uid"`
	TransactionID string    `json:"transaction_id" storage:"transaction_id"`
	RequestID     string    `json:"request_id" storage:"request_id"`
	Currency      string    `json:"currency" storage:"currency"`
	Provider      string    `json:"provider" storage:"provider"`
	Amount        int       `json:"amount" storage:"amount"`
	PaymentDt     time.Time `json:"payment_dt" storage:"payment_dt"`
	Bank          string    `json:"bank" storage:"bank"`
	DeliveryCost  int       `json:"delivery_cost" storage:"delivery_cost"`
	GoodsTotal    int       `json:"goods_total" storage:"goods_total"`
	CustomFee     int       `json:"custom_fee" storage:"custom_fee"`
}

// Item represents the items table
type Item struct {
	ItemID      int    `json:"item_id" storage:"item_id"`
	OrderUid    string `json:"order_uid" storage:"order_uid"`
	ChrtID      int    `json:"chrt_id" storage:"chrt_id"`
	TrackNumber string `json:"track_number" storage:"track_number"`
	Price       int    `json:"price" storage:"price"`
	RID         string `json:"rid" storage:"rid"`
	Name        string `json:"name" storage:"name"`
	Sale        int    `json:"sale" storage:"sale"`
	Size        string `json:"size" storage:"size"`
	TotalPrice  int    `json:"total_price" storage:"total_price"`
	NMID        int    `json:"nm_id" storage:"nm_id"`
	Brand       string `json:"brand" storage:"brand"`
	Status      int    `json:"status" storage:"status"`
}
