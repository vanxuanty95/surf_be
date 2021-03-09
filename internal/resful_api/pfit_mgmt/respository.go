package pfit_mgmt

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"surf_be/internal/app/bot"
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

func (rp *Repository) GetBotByEmailAndPair(ctx context.Context, email, pair string) (*bot.Bot, error) {
	botResult := bot.Bot{}
	filter := bson.M{"email": email, "pair": pair}
	if err := rp.Collection.FindOne(ctx, filter).Decode(&botResult); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &botResult, nil
}

func (rp *Repository) InsertBot(ctx context.Context, bot bot.Bot) error {
	_, err := rp.Collection.InsertOne(ctx, bot)
	if err != nil {
		log.Printf("expect to insert bot to database, but got error: %v", err)
		return err
	}
	return nil
}

func (rp *Repository) GetBotDynamicFields(ctx context.Context, fields map[string]interface{}) ([]*bot.Bot, error) {
	var botResult []*bot.Bot

	filter := make(bson.M, 1)
	for field, value := range fields {
		filter[field] = value
	}

	cursor, err := rp.Collection.Find(ctx, filter)
	if err != nil {
		return botResult, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &botResult); err != nil {
		return botResult, err
	}

	return botResult, nil
}
