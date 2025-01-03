package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/IlyasAtabaev731/L0/internal/app/apiserver"
	"github.com/IlyasAtabaev731/L0/internal/cache"
	"github.com/IlyasAtabaev731/L0/internal/config"
	"github.com/IlyasAtabaev731/L0/internal/kafka"
	"github.com/bxcodec/faker/v3"
	_ "github.com/lib/pq"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting application",
		slog.String("env", cfg.Env),
		slog.Any("cfg", cfg),
		slog.Int("port", cfg.ApiPort),
	)

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Pass,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Db,
	)

	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Error("Failed to connect to database: %v", err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Error("Failed to close database connection: %v", err)
		}
	}(db)
	var inMemoryCache = &sync.Map{}
	// Загрузка кэша
	err = cache.LoadCacheFromDB(db, inMemoryCache, log)
	if err != nil {
		log.Error("Failed to load cache from database: %v", err)
	}
	log.Info("Cache loaded successfully")

	// Подключение к Kafka
	brokers := []string{"localhost:9092"}
	topic := "orders"

	go fakeProduce(brokers, topic, log)

	log.Info("Starting Kafka consumer...")
	go kafka.ConsumeKafkaMessages(brokers, topic, db, inMemoryCache, log)

	s := apiserver.New(cfg, log, inMemoryCache, db)

	if err := s.Start(); err != nil {
		log.Error("Failed to start application", err)
	}

	//	Graceful shutdown

	//stop := make(chan os.Signal, 1)
	//signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	//
	//call := <-stop
	//
	//err1 := s.Stop()
	//if err1 != nil {
	//	log.Error("Failed to stop application", err1)
	//}
	//
	//log.Info("Received", slog.String("signal", call.String()))
	//log.Info("Application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func generateFakeOrder() cache.Order {
	var delivery cache.Delivery
	var payment cache.Payment
	var items []cache.Item

	// Используем faker для генерации данных
	delivery = cache.Delivery{
		Name:    faker.Name(),
		Phone:   faker.E164PhoneNumber(),
		Zip:     strconv.Itoa(rand.Intn(100000)),
		City:    faker.Word(),
		Address: faker.Word(),
		Region:  faker.Word(),
		Email:   faker.Email(),
	}

	payment = cache.Payment{
		TransactionID: strconv.Itoa(rand.Intn(100000)),
		Amount:        rand.Intn(1000) + 1,
		RequestID:     strconv.Itoa(rand.Intn(100000)),
		Currency:      "RUB",
		Provider:      "wbpay",
		PaymentDT:     int64(rand.Intn(100000000)),
		Bank:          "alpha",
		DeliveryCost:  rand.Intn(1000) + 1,
		GoodsTotal:    rand.Intn(1000) + 1,
		CustomFee:     0,
	}

	// Создаем 1-3 фейковых товара
	for i := 0; i < rand.Intn(3)+1; i++ {
		var item cache.Item
		faker.FakeData(&item)
		item.ChrtID = rand.Intn(100000)
		item.Price = rand.Intn(1000) + 1
		item.TotalPrice = item.Price - rand.Intn(100)
		items = append(items, item)
	}

	return cache.Order{
		OrderUID:          faker.UUIDHyphenated(),
		TrackNumber:       faker.Word(),
		Entry:             "WBIL",
		Delivery:          delivery,
		Payment:           payment,
		Items:             items,
		Locale:            "en",
		CustomerID:        faker.UUIDHyphenated(),
		DeliveryService:   "meest",
		DateCreated:       time.Now().Format(time.RFC3339),
		InternalSignature: faker.UUIDHyphenated(),
		SmID:              rand.Intn(100),
		OofShard:          strconv.Itoa(rand.Intn(10)),
		ShardKey:          strconv.Itoa(rand.Intn(10)),
	}
}

func fakeProduce(brokers []string, topic string, log *slog.Logger) {
	producer, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		log.Error("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	log.Info("Starting to produce fake orders...")

	for {
		order := generateFakeOrder()

		// Преобразуем заказ в JSON
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Error("Failed to marshal order: %v", err)
			continue
		}

		// Отправляем сообщение в Kafka
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(orderJSON),
		}

		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Error("Failed to send message to Kafka: %v", err)
			continue
		}

		log.Info("Order %s sent to partition %d at offset %d", order.OrderUID, partition, offset)

		// Задержка перед следующим сообщением
		time.Sleep(10 * time.Second)
	}
}
