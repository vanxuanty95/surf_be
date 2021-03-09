package binance

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"surf_be/internal/configuration"
)

type (
	Repository struct {
		Config     configuration.Config
		Collection *mongo.Collection
	}
)

func NewRepository(config configuration.Config, client mongo.Client) Repository {
	surfCollection := client.Database(config.Server.DataBase.Mongo.Database).Collection(config.Server.DataBase.Mongo.Collection.Bot)
	return Repository{
		Config:     config,
		Collection: surfCollection,
	}
}

func (rp *Repository) DeleteBotByID(ctx context.Context, id string) error {
	filter := bson.M{"id": id}
	_, err := rp.Collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("expect to insert remove bot to database, but got error: %v", err)
		return err
	}
	return nil
}
