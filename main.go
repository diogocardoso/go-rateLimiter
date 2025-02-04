package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/diogocardoso/rate-limiter/config"
	"github.com/diogocardoso/rate-limiter/logic"
	"github.com/diogocardoso/rate-limiter/middleware"
	"github.com/diogocardoso/rate-limiter/storage"
)

func main() {
	cfg := config.LoadConfig()

	redisStorage, err := storage.NewRedisStorage(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	// Definindo limite baixo para testes (por exemplo, 2 requisições por segundo)
	limiter := logic.NewRateLimiter(redisStorage, 2, 2, int64(cfg.BlockTime))

	handler := middleware.RateLimiterMiddleware(limiter, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Request allowed"))
	}))

	// Adicionar handler para a rota específica
	mux := http.NewServeMux()
	mux.Handle("/api", handler)

	port := cfg.Port
	log.Printf("Server running on port %s...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		log.Fatal(err)
	}
}
