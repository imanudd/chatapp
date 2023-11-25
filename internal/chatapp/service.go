package chatapp

import (
	"context"
	"finalproject/config"

	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	cfg  *config.Config
	db   *mongo.Database
	repo IRepository
}

func NewService(cfg *config.Config, db *mongo.Database, repo IRepository) Service {
	return Service{cfg, db, repo}
}

func (s Service) SendMessage(ctx context.Context) (*Users, error) {

	return s.repo.SendMessage(ctx)
}
