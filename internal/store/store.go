package store

import (
	"github.com/Calavrat/TestMedods/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	GetUser(id string) (model.User, error)
	AddRefresh(model.User) error
}
type Repository struct {
	Authorization
}

func NewRepository(db *mongo.Database, collection string) *Repository {
	return &Repository{
		Authorization: NewAuthorization(db, collection),
	}
}
