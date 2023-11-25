package chatapp

import (
	"context"
	"finalproject/config"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	db  *mongo.Database
	cfg *config.Config
}

func NewRepository(db *mongo.Database, cfg *config.Config) Repository {
	return Repository{db, cfg}
}

func (r Repository) SendMessage(ctx context.Context) (*Users, error) {

	return nil, nil
}
