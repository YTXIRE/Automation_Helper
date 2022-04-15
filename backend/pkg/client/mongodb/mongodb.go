package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Ctx      context.Context
	Host     string
	Port     string
	Username string
	Password string
	Database string
	AuthDB   string
}

func NewClient(config *Config) (db *mongo.Database, err error) {
	var mongoDBURL string
	var isAuth bool
	if config.Username == "" && config.Password == "" {
		mongoDBURL = fmt.Sprintf("mongodb://%s:%s", config.Host, config.Port)
	} else {
		isAuth = true
		mongoDBURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", config.Username, config.Password, config.Host, config.Port)
	}

	clientOptions := options.Client().ApplyURI(mongoDBURL)

	if isAuth {
		if config.AuthDB == "" {
			config.AuthDB = config.Database
		}
		clientOptions.SetAuth(options.Credential{
			AuthSource: config.AuthDB,
			Username:   config.Username,
			Password:   config.Password,
		})
	}

	client, err := mongo.Connect(config.Ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongoDB due to error: %v", err)
	}

	if err := client.Ping(config.Ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping to mongoDB due to error: %v", err)
	}

	return client.Database(config.Database), nil
}
