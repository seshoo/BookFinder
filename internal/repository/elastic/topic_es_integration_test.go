package elastic_test

import (
	"context"
	"github.com/seshoo/bookFinder/internal/domain"
	"github.com/seshoo/bookFinder/internal/repository/elastic"
	"github.com/seshoo/bookFinder/pkg/elasticsearch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var mappingString = `{
  "mappings": {
    "properties": {
      "id": {
        "type": "keyword"
      },
      "title": {
        "type": "text",
        "analyzer": "russian"
      },
      "link": {
        "type": "keyword"
      },
      "text": {
        "type": "text",
        "analyzer": "russian"
      }
    }
  }
}`

func TestUserRepository(t *testing.T) {
	var (
		esUsername = "test"
		esPassword = "test"
		index      = "test_topic_index"
	)

	esHost, teardown := startTestContainer(t, esUsername, esPassword)
	defer teardown()

	esClient := getClient(t, esHost, esUsername, esPassword)

	createIndex(t, esClient, index)
	setMapping(t, esClient, index, mappingString)

	db, err := elasticsearch.NewElastic(elasticsearch.Config{
		Host:      esHost,
		Username:  esUsername,
		Password:  esPassword,
		IndexName: index,
	})
	if err != nil {
		t.Fatalf("failed to create test database: %s", err)
	}

	topicRepository := elastic.NewTopicRepository(db)

	t.Run("GetById return topic if it exists", func(t *testing.T) {

		expectedTopic := &domain.Topic{
			Id:    "1",
			Title: "test title",
			Link:  "test link",
			Text:  "test text",
		}

		_, err := esClient.Index(index).
			Id(expectedTopic.Id).
			Request(expectedTopic).
			Do(context.TODO())
		if err != nil {
			t.Fatalf("failed to index document: %s", err)
		}

		actualTopic, err := topicRepository.GetById(expectedTopic.Id)
		require.NoError(t, err)

		assert.Equal(t, expectedTopic, actualTopic)

	})
}
