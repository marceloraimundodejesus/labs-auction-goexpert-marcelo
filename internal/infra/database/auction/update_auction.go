package auction

import (
	"context"

	"fullcycle-auction_go/internal/entity/auction_entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ar *AuctionRepository) UpdateAuctionStatus(ctx context.Context, id string, status auction_entity.AuctionStatus) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = ar.Collection.UpdateByID(ctx, oid, bson.M{
		"$set": bson.M{"status": status},
	})
	return err
}
