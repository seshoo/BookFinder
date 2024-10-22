package elasticsearch

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/putmapping"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type Config struct {
	Host      string
	Username  string
	Password  string
	IndexName string
}

type Elastic struct {
	indexName string
	client    *elasticsearch.TypedClient
}

func NewElastic(config Config) (*Elastic, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			config.Host,
		},
		Username: config.Username,
		Password: config.Password,
	}

	es, err := elasticsearch.NewTypedClient(cfg)

	if err != nil {
		return nil, err
	}

	return &Elastic{client: es, indexName: config.IndexName}, nil
}

func (el *Elastic) Prepare(ctx context.Context, mappingByte []byte) error {
	var (
		err     error
		mapping types.TypeMapping
	)

	err = json.Unmarshal(mappingByte, &mapping)
	if err != nil {
		return err
	}

	_, err = el.client.Indices.Get(el.indexName).Do(ctx)
	if err != nil {
		err = el.createIndexWithMapping(ctx, &mapping)
	} else {
		err = el.updateMapping(ctx, &mapping)
	}

	return nil
}

func (el *Elastic) Get(ctx context.Context, id string, document interface{}) error {
	res, err := el.client.Get(el.indexName, id).Do(ctx)
	if err != nil {
		return err
	}

	err = json.Unmarshal(res.Source_, &document)
	_ = res
	return err
}

func (el *Elastic) createIndexWithMapping(ctx context.Context, mp *types.TypeMapping) error {
	_, err := el.client.Indices.Create(el.indexName).Request(&create.Request{
		Mappings: mp,
	}).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (el *Elastic) updateMapping(ctx context.Context, mp *types.TypeMapping) error {
	_, err := el.client.Indices.PutMapping(el.indexName).Request(&putmapping.Request{
		Properties: mp.Properties,
	}).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
