package apiserver

import (
	"database/sql"
	"encoding/json"
	"github.com/IlyasAtabaev731/L0/internal/config"
	"github.com/gorilla/mux"
	"io"
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
	db     *sql.DB
}

func New(config *config.Config, logger *slog.Logger, cache *sync.Map, db *sql.DB) *APIServer {
	return &APIServer{
		config: config,
		logger: logger,
		cache:  cache,
		router: mux.NewRouter(),
		db:     db,
	}
}

func (s *APIServer) Start() error {
	s.logger.Info("Starting server", slog.String("port", strconv.Itoa(s.config.ApiPort)))

	s.configureRouter()

	return http.ListenAndServe(s.config.ApiHost+":"+strconv.Itoa(s.config.ApiPort), s.router)
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/orderDetails", s.handleOrderDetails())
	s.router.HandleFunc("/order/{orderUid}", s.orderHandler).Methods("GET")
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

func (s *APIServer) orderHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем orderUID из URL
	vars := mux.Vars(r)
	orderUID := vars["orderUid"]

	// Ищем в кэше
	order, found := s.cache.Load(orderUID)
	if !found {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Преобразуем данные в JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		http.Error(w, "Failed to marshal order data", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err1 := w.Write(orderJSON)
	if err1 != nil {
		s.logger.Error("Failed to write response", err1)
	}
}
