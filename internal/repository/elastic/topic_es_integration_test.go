package elastic_test

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
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
		err        error
		esUsername = "test"
		esPassword = "test"
		indexName  = "test_topic_index"
	)

	esHost, teardown := startTestContainer(t, esUsername, esPassword)
	defer teardown()

	esClient := getClient(t, esHost, esUsername, esPassword)

	createIndex(t, esClient, indexName)
	//setMapping(t, esClient, indexName, mappingString)

	db, err := elasticsearch.NewElastic(elasticsearch.Config{
		Host:      esHost,
		Username:  esUsername,
		Password:  esPassword,
		IndexName: indexName,
	})
	if err != nil {
		t.Fatalf("failed to create test database: %s", err)
	}

	topicRepository := elastic.NewTopicRepository(db)

	t.Run("GetById return topic if it exists", func(t *testing.T) {
		removeAllFromIndex(t, esClient, indexName)

		expectedTopic := &domain.Topic{
			Id:    "1",
			Title: "test title",
			Link:  "test link",
			Text:  "test text",
		}

		_, err := esClient.Index(indexName).
			Id(expectedTopic.Id).
			Request(expectedTopic).
			Do(context.TODO())
		if err != nil {
			t.Fatalf("failed to index document: %s", err)
		}
		defer removeAllFromIndex(t, esClient, indexName)

		actualTopic, err := topicRepository.GetById(expectedTopic.Id)

		assert.Equal(t, expectedTopic, actualTopic)
		require.NoError(t, err)
	})

	t.Run("GetById return topic if it not exists", func(t *testing.T) {
		removeAllFromIndex(t, esClient, indexName)

		actualTopic, err := topicRepository.GetById("2")

		assert.Nil(t, actualTopic)
		require.Error(t, err)
	})

	t.Run("GetList return topics", func(t *testing.T) {
		removeAllFromIndex(t, esClient, indexName)

		bulk := esClient.Bulk().Index(indexName).Refresh(refresh.True)

		doc1 := &domain.Topic{
			Id:    "1",
			Title: "test title 1",
			Link:  "test link 1",
			Text:  "test text 1",
		}

		doc2 := &domain.Topic{
			Id:    "2",
			Title: "test title 2",
			Link:  "test link 2",
			Text:  "test text 2",
		}

		err = bulk.CreateOp(types.CreateOperation{Id_: &doc1.Id}, doc1)
		if err != nil {
			t.Fatalf("failed to create bulk operation: %s", err)
		}

		err = bulk.CreateOp(types.CreateOperation{Id_: &doc2.Id}, doc2)
		if err != nil {
			t.Fatalf("failed to create bulk operation: %s", err)
		}

		_, err = bulk.Do(context.Background())

		if err != nil {
			t.Fatalf("failed to bulk index documents: %s", err)
		}
	})
}
