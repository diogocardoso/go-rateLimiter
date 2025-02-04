package logic_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/diogocardoso/rate-limiter/config"
	"github.com/diogocardoso/rate-limiter/logic"
	"github.com/diogocardoso/rate-limiter/storage"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiterByIP(t *testing.T) {
	// Criar configurações de teste
	cfg := config.Config{
		RedisAddr:      "localhost:6379",
		RedisPassword:  "1234",
		RedisDB:        0,
		RateLimitIP:    2,
		RateLimitToken: 3,
		BlockTime:      60, // int
	}

	// Criar o storage com Redis diretamente passando os parâmetros
	redisStorage, err := storage.NewRedisStorage(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB) // Passando os parâmetros corretos
	if err != nil {
		t.Fatalf("Erro ao conectar ao Redis: %v", err)
	}

	// Criar o rate limiter passando BlockTime como int64
	limiter := logic.NewRateLimiter(redisStorage, cfg.RateLimitIP, cfg.RateLimitToken, int64(cfg.BlockTime)) // Convertendo para int64

	// Agora você pode prosseguir com o teste
	// Exemplo: testar se o limite de requisições por IP está funcionando
	t.Run("TestRateLimiterByIP", func(t *testing.T) {
		// Teste para IP
		// Use o limiter aqui para validar o comportamento
		assert.NotNil(t, limiter)
	})
}

func TestRateLimiterByToken(t *testing.T) {
	// Criar configurações de teste
	cfg := config.Config{
		RedisAddr:      "localhost:6379",
		RedisPassword:  "1234",
		RedisDB:        0,
		RateLimitIP:    10, // Limite de requisições por IP
		RateLimitToken: 3,  // Alterado de 20 para 3 requisições
		BlockTime:      60, // Tempo de bloqueio em segundos
	}

	// Criar o storage com Redis diretamente passando os parâmetros
	redisStorage, err := storage.NewRedisStorage(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		t.Fatalf("Erro ao conectar ao Redis: %v", err)
	}

	// Criar o rate limiter passando BlockTime como int64
	limiter := logic.NewRateLimiter(redisStorage, cfg.RateLimitIP, cfg.RateLimitToken, int64(cfg.BlockTime))

	// Criar um servidor HTTP simulando as requisições
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificando o token no cabeçalho
		token := r.Header.Get("API_KEY")
		if token == "" {
			http.Error(w, "Token não fornecido", http.StatusUnauthorized)
			return
		}

		// Verifica o limite do token
		allowed, err := limiter.IsAllowed(r)
		if err != nil || !allowed {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Se passou na verificação do rate limiter, resposta com sucesso
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Requisição permitida"))
	})

	// Teste 1: Realizando 3 requisições com o mesmo token (deve passar)
	req, _ := http.NewRequest("GET", "/api", nil)
	req.Header.Set("API_KEY", "token123")

	// Enviar a primeira requisição
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Enviar a segunda requisição
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Enviar a terceira requisição
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Teste 2: Realizando a quarta requisição (deve ser bloqueada)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)

	// Teste 3: Realizando uma requisição com um token diferente (deve passar)
	req.Header.Set("API_KEY", "token456")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Após cada requisição
	time.Sleep(100 * time.Millisecond) // Pequeno delay para garantir ordem das requisições
}
