package service

import (
	"github.com/Calavrat/TestMedods/internal/store"
)

type Authorization interface {
	GenerateToken(userid string) (*tokenDetails, error)
	ParseTokens(ts *map[string]string) (string, error)
}

type Service struct {
	Authorization
}

func NewService(repository *store.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repository),
	}
}
