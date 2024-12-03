package cache

import (
	"database/sql"
	"log"
	"sync"
)

// Структура заказа
type Order struct {
	OrderUID        string   `json:"order_uid"`
	TrackNumber     string   `json:"track_number"`
	Entry           string   `json:"entry"`
	Delivery        Delivery `json:"delivery"`
	Payment         Payment  `json:"payment"`
	Items           []Item   `json:"items"`
	Locale          string   `json:"locale"`
	CustomerID      string   `json:"customer_id"`
	DeliveryService string   `json:"delivery_service"`
	DateCreated     string   `json:"date_created"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	TransactionID string `json:"transaction"`
	Currency      string `json:"currency"`
	Provider      string `json:"provider"`
	Amount        int    `json:"amount"`
	PaymentDT     int64  `json:"payment_dt"`
	Bank          string `json:"bank"`
	DeliveryCost  int    `json:"delivery_cost"`
	GoodsTotal    int    `json:"goods_total"`
	CustomFee     int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NMID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

// Загрузка данных из базы данных
func LoadCacheFromDB(db *sql.DB, cache *sync.Map) error {
	rows, err := db.Query("SELECT * FROM orders")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		// Здесь извлекаются данные из таблицы orders, delivery и payment
		// и заполняется структура Order
		// Пример:
		err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.CustomerID, &order.DeliveryService, &order.DateCreated)
		if err != nil {
			log.Println("Error scanning order:", err)
			continue
		}

		// Загрузка доставки
		err = db.QueryRow("SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1", order.OrderUID).
			Scan(&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
		if err != nil {
			log.Println("Error loading delivery:", err)
			continue
		}

		// Загрузка оплаты
		err = db.QueryRow("SELECT transaction, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid = $1", order.OrderUID).
			Scan(&order.Payment.TransactionID, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDT, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee)
		if err != nil {
			log.Println("Error loading payment:", err)
			continue
		}

		// Загрузка товаров
		itemsRows, err := db.Query("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid = $1", order.OrderUID)
		if err != nil {
			log.Println("Error loading items:", err)
			continue
		}
		for itemsRows.Next() {
			var item Item
			err = itemsRows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NMID, &item.Brand, &item.Status)
			if err != nil {
				log.Println("Error scanning item:", err)
				continue
			}
			order.Items = append(order.Items, item)
		}
		itemsRows.Close()

		// Сохранение в кэш
		cache.Store(order.OrderUID, order)
	}

	return nil
}
