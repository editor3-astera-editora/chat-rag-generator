package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/editor3-astera-editora/chat-rag/internal/service"
)

type ChatHandler struct {
	Service *service.ChatService
}

func NewChatHandler(service *service.ChatService) *ChatHandler {
	return &ChatHandler{Service: service}
}

func (h *ChatHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID  string `json:"user_id"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(" Erro ao decodificar JSON", err)
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	resp, sources, err := h.Service.ProcessMessage(r.Context(), req.UserID, req.Message)
	if err != nil {
		log.Println(" Erro interno:", err)
		http.Error(w, "Erro interno: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(" Resposta gerada:", resp)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"response": resp,
		"sources":  sources,
	})
}
