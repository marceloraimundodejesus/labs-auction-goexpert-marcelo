//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/auction"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAutoCloseAuction(t *testing.T) {
	t.Setenv("AUCTION_DURATION", "2s")
	t.Setenv("AUCTION_INTERVAL", "500ms")

	if os.Getenv("MONGODB_URL") == "" {
		t.Setenv("MONGODB_URL", "mongodb://admin:admin@127.0.0.1:27017/auctions?authSource=admin")
	}

	ctx := context.Background()

	db, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		t.Fatalf("falha conectando no Mongo: %v", err)
	}

	// Se o nome do DB vier vazio, forçamos "auctions"
	if db.Name() == "" {
		db = db.Client().Database("auctions")
	}

	// 1) Ping (sanity check)
	if err := db.Client().Ping(ctx, nil); err != nil {
		t.Fatalf("ping ao Mongo falhou: %v", err)
	}

	// 2) Insert cru para capturar erros reais do Mongo
	rawID := "it-" + uuid.NewString()
	rawDoc := bson.M{"_id": rawID, "smoke": true, "ts": time.Now().Unix()}
	coll := db.Collection("auctions")
	if _, err := coll.InsertOne(ctx, rawDoc); err != nil {
		t.Fatalf("mongo insert raw falhou: %v", err)
	}
	_, _ = coll.DeleteOne(ctx, bson.M{"_id": rawID})

	// 3) Fluxo via repositório
	repo := auction.NewAuctionRepository(db)

	auctionEnt, ierr := auction_entity.CreateAuction(
		"Integration Test - AutoClose",
		"electronics",
		"desc 1234567890",
		auction_entity.Used,
	)
	if ierr != nil {
		t.Fatalf("erro criando entidade de auction: %v", ierr)
	}

	if err := repo.CreateAuction(ctx, auctionEnt); err != nil {
		t.Fatalf("erro inserindo auction via repo: %v", err)
	}

	// Espera > duração p/ fechamento automático
	time.Sleep(3 * time.Second)

	got, ierr := repo.FindAuctionById(ctx, auctionEnt.Id)
	if ierr != nil {
		t.Fatalf("erro buscando auction por id: %v", ierr)
	}

	if got.Status != auction_entity.Completed {
		t.Fatalf("status esperado %v, obtido %v", auction_entity.Completed, got.Status)
	}
}
