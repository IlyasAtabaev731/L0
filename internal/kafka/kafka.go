package kafka

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/IBM/sarama"
	"log/slog"
	"sync"

	"github.com/IlyasAtabaev731/L0/internal/cache"
)

func ConsumeKafkaMessages(brokers []string, topic string, db *sql.DB, inMemoryCache *sync.Map, log *slog.Logger) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	consumerGroup, err := sarama.NewConsumerGroup(brokers, "orders_consumer_group", config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer group: %v", err)
	}
	defer consumerGroup.Close()

	handler := &consumerHandler{
		db:            db,
		inMemoryCache: inMemoryCache,
	}

	for {
		err := consumerGroup.Consume(context.Background(), []string{topic}, handler)
		if err != nil {
			log.Printf("Error in Kafka consumer: %v", err)
		}
	}
}

// Custom handler
type consumerHandler struct {
	db            *sql.DB
	inMemoryCache *sync.Map
}

func (h *consumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var order cache.Order
		err := json.Unmarshal(message.Value, &order)
		if err != nil {
			log.Printf("Failed to parse Kafka message: %v", err)
			continue
		}

		// Сохранение в базу данных
		err = saveOrderToDB(h.db, order)
		if err != nil {
			log.Printf("Failed to save order to database: %v", err)
			continue
		}

		// Сохранение в кэш
		h.inMemoryCache.Store(order.OrderUID, order)

		log.Printf("Processed order %s", order.OrderUID)
		session.MarkMessage(message, "")
	}
	return nil
}

func saveOrderToDB(db *sql.DB, order cache.Order) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Сохранение заказа
	_, err = tx.Exec(
		`INSERT INTO orders (order_uid, track_number, entry, locale, customer_id, internal_signature, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO NOTHING`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.CustomerID, order.InternalSignature, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Сохранение доставки
	_, err = tx.Exec(
		`INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (order_uid) DO NOTHING`,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City,
		order.Delivery.Address, order.Delivery.Region, order.Delivery.Email,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Сохранение оплаты
	_, err = tx.Exec(
		`INSERT INTO payments (order_uid, transaction_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (order_uid) DO NOTHING`,
		order.OrderUID, order.Payment.TransactionID, order.Payment.Currency, order.Payment.Provider,
		order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Сохранение товаров
	for _, item := range order.Items {
		_, err = tx.Exec(
			`INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			ON CONFLICT DO NOTHING`,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale,
			item.Size, item.TotalPrice, item.NMID, item.Brand, item.Status,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func subscribeToKafka(brokers []string, topic string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	return sarama.NewConsumerGroup(brokers, "group-id", config)
}

func publishToKafka(brokers []string, topic string, message []byte) error {
	producer, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		return err
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}
	_, _, err = producer.SendMessage(msg)
	return err
}
