package handlers

import (
	"kasir-api/services"
	"kasir-api/models"
	"net/http"
	"encoding/json"

)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// HandleTransaction - Get /api/checkout
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Create(w, r)		
	default:
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request Body", http.StatusBadRequest)
		return
	}
	transaction, err := h.service.Create(req.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}