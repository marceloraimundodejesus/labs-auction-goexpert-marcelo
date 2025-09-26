package config

import (
	"os"
	"testing"
	"time"
)

func TestAuctionDuration_Default(t *testing.T) {
	_ = os.Unsetenv("AUCTION_DURATION")
	if got := AuctionDuration(); got != 5*time.Minute {
		t.Fatalf("default mismatch: got %v want %v", got, 5*time.Minute)
	}
}

func TestAuctionDuration_CustomValid(t *testing.T) {
	t.Setenv("AUCTION_DURATION", "90s")
	if got := AuctionDuration(); got != 90*time.Second {
		t.Fatalf("valid env mismatch: got %v want %v", got, 90*time.Second)
	}
}

func TestAuctionDuration_InvalidFallsBack(t *testing.T) {
	t.Setenv("AUCTION_DURATION", "xpto")
	if got := AuctionDuration(); got != 5*time.Minute {
		t.Fatalf("invalid should fallback: got %v want %v", got, 5*time.Minute)
	}
}

func TestAuctionInterval_Default(t *testing.T) {
	_ = os.Unsetenv("AUCTION_INTERVAL")
	if got := AuctionInterval(); got != 20*time.Second {
		t.Fatalf("default mismatch: got %v want %v", got, 20*time.Second)
	}
}

func TestAuctionInterval_CustomValid(t *testing.T) {
	t.Setenv("AUCTION_INTERVAL", "1s")
	if got := AuctionInterval(); got != 1*time.Second {
		t.Fatalf("valid env mismatch: got %v want %v", got, 1*time.Second)
	}
}
