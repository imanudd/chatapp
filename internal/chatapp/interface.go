package chatapp

import (
	"context"
)

type IService interface {
	SendMessage(ctx context.Context) (*Users, error)
	// ReceipeMessage(ctx context.Context) error
}

type IRepository interface {
	SendMessage(ctx context.Context) (*Users, error)
}
