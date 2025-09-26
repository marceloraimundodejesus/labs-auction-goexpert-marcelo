package auction

import (
	"context"

	"fullcycle-auction_go/internal/entity/auction_entity"

	"go.mongodb.org/mongo-driver/bson"
)

func (ar *AuctionRepository) UpdateAuctionStatus(ctx context.Context, id string, status auction_entity.AuctionStatus) error {
	_, err := ar.Collection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"status": status}},
	)
	return err
}
