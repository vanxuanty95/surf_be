package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"surf_be/internal/configuration"
	"time"
)

type Mongo struct {
	Config configuration.Config
	Client *mongo.Client
}

func (mg *Mongo) Init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	connectionOption := options.ClientOptions{
		Hosts: mg.Config.Server.DataBase.Mongo.Host,
		Auth: &options.Credential{
			Username:      mg.Config.Server.DataBase.Mongo.Username,
			Password:      mg.Config.Server.DataBase.Mongo.Password,
			AuthMechanism: "SCRAM-SHA-1",
			AuthSource:    mg.Config.Server.DataBase.Mongo.AuthDatabase,
		},
	}
	client, err := mongo.Connect(ctx, &connectionOption)
	if err != nil {
		cancel()
		log.Fatal(err)
	}
	mg.Client = client
}

func (mg *Mongo) Disconnect() {
	if err := mg.Client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}
