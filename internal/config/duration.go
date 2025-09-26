package config

import (
	"os"
	"time"
)

// getDuration lê uma env como time.Duration; se vazia/ inválida, usa default.
func getDuration(env string, def time.Duration) time.Duration {
	val := os.Getenv(env)
	if val == "" {
		return def
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return def
	}
	return d
}

// Tempo total de vida do leilão (ex.: 5m). Default: 5 minutos.
func AuctionDuration() time.Duration {
	return getDuration("AUCTION_DURATION", 5*time.Minute)
}

// Intervalo do ticker que verifica leilões vencidos. Default: 20s.
func AuctionInterval() time.Duration {
	return getDuration("AUCTION_INTERVAL", 20*time.Second)
}
