package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"tg-tictactoe-miniapp/internal/bot"
	"tg-tictactoe-miniapp/internal/webapp"
)

func main() {
	// Load .env for local development convenience; ignore errors in production environments.
	_ = godotenv.Load()

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}
	log.Printf("TELEGRAM_BOT_TOKEN=%s", token)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("PORT=%s", port)

	// Initialize Telegram bot client.
	botClient := bot.NewBotClient(token)

	// Initialize Gin router.
	router := gin.Default()

	// Register web and API routes.
	webapp.RegisterRoutes(router, botClient)

	// Start HTTP server.
	go func() {
		addr := ":" + port
		log.Printf("Starting HTTP server at %s...", addr)
		if err := router.Run(addr); err != nil {
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	// Wait for termination signals.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	log.Printf("Received signal %s, shutting down...", sig)
}
