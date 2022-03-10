package store

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	AuthDB   string
}

func NewMongoDB(ctx context.Context, cfg *Config) (*mongo.Database, error) {
	var mongoDBURL string
	var isAuth bool
	if cfg.Username == "" && cfg.Password == "" {
		mongoDBURL = fmt.Sprintf("mongodb://%s:%s", cfg.Host, cfg.Port)
	} else {
		isAuth = true
		mongoDBURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	}

	clientOptions := options.Client().ApplyURI(mongoDBURL)

	if isAuth {
		if cfg.AuthDB == "" {
			cfg.AuthDB = cfg.Database
		}
		clientOptions.SetAuth(options.Credential{
			AuthSource: cfg.AuthDB,
			Username:   cfg.Username,
			Password:   cfg.Password,
		})
	}

	//Connect
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB due to error: %v", err)
	}

	//Ping
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB due to error: %v", err)
	}

	return client.Database(cfg.Database), nil
}
