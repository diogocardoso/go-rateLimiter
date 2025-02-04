package logic

import (
	"net/http"
	"strings"

	"github.com/diogocardoso/rate-limiter/storage"
)

type RateLimiter struct {
	storage       storage.Storage
	ipLimit       int
	tokenLimit    int
	blockDuration int64
}

func NewRateLimiter(storage storage.Storage, ipLimit, tokenLimit int, blockDuration int64) *RateLimiter {
	return &RateLimiter{
		storage:       storage,
		ipLimit:       ipLimit,
		tokenLimit:    tokenLimit,
		blockDuration: blockDuration,
	}
}

func (rl *RateLimiter) IsAllowed(r *http.Request) (bool, error) {
	// Verificar IP
	ip := strings.Split(r.RemoteAddr, ":")[0]

	// Verificar token
	token := r.Header.Get("API_KEY")

	// Verificar se IP ou token estão bloqueados
	isIPBlocked, err := rl.storage.IsBlocked(ip)
	if err != nil {
		return false, err
	}

	if token != "" {
		isTokenBlocked, err := rl.storage.IsBlocked(token)
		if err != nil {
			return false, err
		}
		if isTokenBlocked {
			return false, nil
		}
	}

	if isIPBlocked {
		return false, nil
	}

	// Incrementar contadores
	ipCount, err := rl.storage.Increment(ip)
	if err != nil {
		return false, err
	}

	if token != "" {
		tokenCount, err := rl.storage.Increment(token)
		if err != nil {
			return false, err
		}
		if tokenCount > int64(rl.tokenLimit) {
			err = rl.storage.Block(token, rl.blockDuration)
			if err != nil {
				return false, err
			}
			return false, nil
		}
	}

	if ipCount > int64(rl.ipLimit) {
		err = rl.storage.Block(ip, rl.blockDuration)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}
