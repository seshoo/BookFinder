package repository

import (
	"github.com/seshoo/bookFinder/internal/domain"
)

type Topics interface {
	GetList() ([]domain.Topic, error)
	GetById(id string) (*domain.Topic, error)
}

type Repositories interface {
	Topic() Topics
}
