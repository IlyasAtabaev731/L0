package apiserver

import (
	"database/sql"
	"encoding/json"
	"github.com/IlyasAtabaev731/L0/internal/cache"
	"github.com/IlyasAtabaev731/L0/internal/config"
	"github.com/gorilla/mux"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
)

type APIServer struct {
	config *config.Config
	logger *slog.Logger
	cache  *sync.Map
	router *mux.Router
}

func New(config *config.Config, logger *slog.Logger, cache *sync.Map) *APIServer {
	return &APIServer{
		config: config,
		logger: logger,
		cache:  cache,
		router: mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	s.logger.Info("Starting server", slog.String("port", strconv.Itoa(s.config.ApiPort)))

	s.configureRouter()

	return http.ListenAndServe(s.config.ApiHost+":"+strconv.Itoa(s.config.ApiPort), s.router)
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/orderDetails", s.handleOrderDetails())
	//s.router.HandleFunc("/order/")
}

func (s *APIServer) handleOrderDetails() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "<!DOCTYPE html>\n<html>\n<head>\n    <title>Order Details</title>\n</head>\n<body>\n    <h1>Order Details</h1>\n    <form id=\"orderForm\">\n        <label for=\"orderUid\">Enter Order UID:</label><br>\n        <input type=\"text\" id=\"orderUid\" name=\"orderUid\"><br>\n        <input type=\"submit\" value=\"Submit\">\n    </form>\n    <div id=\"orderData\"></div>\n\n    <script>\n        document.getElementById('orderForm').onsubmit = function(event) {\n            event.preventDefault();\n            const orderUid = document.getElementById('orderUid').value;\n            fetch(`/order/${orderUid}`)\n                .then(response => response.json())\n                .then(data => {\n                    if (data.error) {\n                        document.getElementById('orderData').innerHTML = `<p>${data.error}</p>`;\n                    } else {\n                        document.getElementById('orderData').innerHTML = `<pre>${JSON.stringify(data, null, 2)}</pre>`;\n                    }\n                })\n                .catch(error => {\n                    console.error('Error:', error);\n                    document.getElementById('orderData').innerHTML = `<p>An error occurred.</p>`;\n                });\n        };\n    </script>\n</body>\n</html>")
		if err != nil {
			s.logger.Error("Failed to write response", err)
		}
	}
}

func (s *APIServer) Stop() error {
	return nil
}

func (s *APIServer) getOrderById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if val, ok := s.cache.Load(id); ok {
		err := json.NewEncoder(w).Encode(val)
		if err != nil {
			return
		}
	} else {
		http.Error(w, "Order not found", http.StatusNotFound)
	}
}

func getOrderHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем параметр из URL
		vars := mux.Vars(r)
		orderUid := vars["orderUid"]

		// Получаем данные заказа из БД
		var order cache.Order
		err := db.QueryRow(
			`SELECT order_uid, track_number, entry, locale, customer_id, delivery_service, date_created 
			 FROM orders 
			 WHERE order_uid = $1`,
			orderUid,
		).Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.CustomerID, &order.DeliveryService, &order.DateCreated)

		if err == sql.ErrNoRows {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("Failed to query order: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Устанавливаем заголовки ответа
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Возвращаем данные в формате JSON
		json.NewEncoder(w).Encode(order)
	}
}
