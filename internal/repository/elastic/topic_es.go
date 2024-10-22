package elastic

import (
	"context"
	"github.com/seshoo/bookFinder/internal/domain"
	"github.com/seshoo/bookFinder/pkg/elasticsearch"
)

type TopicRepository struct {
	db *elasticsearch.Elastic
}

func NewTopicRepository(db *elasticsearch.Elastic) *TopicRepository {
	return &TopicRepository{db: db}
}

func (t TopicRepository) GetList() ([]domain.Topic, error) {

	return []domain.Topic{}, nil
}

func (t TopicRepository) GetById(id string) (*domain.Topic, error) {
	var topic domain.Topic
	err := t.db.Get(context.Background(), id, &topic)
	if err != nil {
		return nil, err
	}
	return &topic, nil

}
