package main

import (
	"log"
	"net/http"

	"github.com/editor3-astera-editora/chat-rag/internal/handler"
	"github.com/editor3-astera-editora/chat-rag/internal/middleware"
	"github.com/editor3-astera-editora/chat-rag/internal/repository"
	"github.com/editor3-astera-editora/chat-rag/internal/service"
	"github.com/editor3-astera-editora/chat-rag/pkg/config"
)

func main() {
	db := config.GetDB()
	defer db.Close()

	repo := repository.NewChatRepository(db)
	chatService := service.NewChatService(repo)
	chatHandler := handler.NewChatHandler(chatService)

	mux := http.NewServeMux()
	mux.HandleFunc("/chat", chatHandler.HandleChat)

	server := middleware.LoggingMiddleware(
		middleware.CORSMiddleware(mux),
	)

	log.Println("Servidor inicado em http://localhost:8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatal(err)
	}
}
