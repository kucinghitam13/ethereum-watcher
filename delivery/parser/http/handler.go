package http

import (
	"encoding/json"
	"net/http"
)

type (
	genericResponse struct {
		Data any `json:"data"`
	}
)

func (this *Handler) GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	currentBlock := this.usecase.GetCurrentBlock(r.Context())
	writeJSON(w, http.StatusOK, currentBlock)
}

func (this *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	address := r.FormValue("address")
	if address == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	status := this.usecase.Subscribe(r.Context(), address)
	writeJSON(w, http.StatusOK, status)
}

func (this *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	address := r.FormValue("address")
	if address == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	transactions := this.usecase.GetTransactions(r.Context(), address)
	writeJSON(w, http.StatusOK, transactions)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	b, _ := json.Marshal(genericResponse{
		Data: data,
	})
	w.Write(b)
}
