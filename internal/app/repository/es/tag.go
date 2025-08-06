package es

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/technoshantoms/mccs-alpha-api/internal/app/types"
	"github.com/elastic/go-elasticsearch/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type tag struct {
	c     *elastic.Client
	index string
}

var Tag = &tag{}

func (es *tag) Register(client *elastic.Client) {
	es.c = client
	es.index = "tags"
}

func (es *tag) Create(id primitive.ObjectID, name string) error {
	body := types.TagESRecord{
		TagID: id.Hex(),
		Name:  name,
	}
	_, err := es.c.Index().
		Index(es.index).
		Id(id.Hex()).
		BodyJson(body).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (es *tag) UpdateOffer(id string, name string) error {
	exists, err := es.c.Exists().Index(es.index).Id(id).Do(context.TODO())
	if err != nil {
		return err
	}
	if !exists {
		body := types.TagESRecord{
			TagID:        id,
			Name:         name,
			OfferAddedAt: time.Now(),
		}
		_, err = es.c.Index().
			Index(es.index).
			Id(id).
			BodyJson(body).
			Do(context.Background())
		if err != nil {
			return err
		}
		return nil
	}

	params := map[string]interface{}{
		"offerAddedAt": time.Now(),
	}
	script := elastic.
		NewScript(`
			ctx._source.offerAddedAt = params.offerAddedAt;
		`).
		Params(params)

	_, err = es.c.Update().
		Index(es.index).
		Id(id).
		Script(script).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (es *tag) UpdateWant(id string, name string) error {
	exists, err := es.c.Exists().Index(es.index).Id(id).Do(context.TODO())
	if err != nil {
		return err
	}
	if !exists {
		body := types.TagESRecord{
			TagID:       id,
			Name:        name,
			WantAddedAt: time.Now(),
		}
		_, err = es.c.Index().
			Index(es.index).
			Id(id).
			BodyJson(body).
			Do(context.Background())
		if err != nil {
			return err
		}
		return nil
	}

	params := map[string]interface{}{
		"wantAddedAt": time.Now(),
	}
	script := elastic.
		NewScript(`
			ctx._source.wantAddedAt = params.wantAddedAt;
		`).
		Params(params)

	_, err = es.c.Update().
		Index(es.index).
		Id(id).
		Script(script).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (es *tag) Update(id primitive.ObjectID, update *types.Tag) error {
	params := map[string]interface{}{
		"name": update.Name,
	}
	script := elastic.
		NewScript(`
			ctx._source.name = params.name;
		`).
		Params(params)

	_, err := es.c.Update().
		Index(es.index).
		Id(id.Hex()).
		Script(script).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (es *tag) DeleteByID(id string) error {
	_, err := es.c.Delete().
		Index(es.index).
		Id(id).
		Do(context.Background())
	if err != nil {
		if elastic.IsNotFound(err) {
			return errors.New("Tag does not exist.")
		}
		return err
	}
	return nil
}

// MatchOffer matches wants for the given offer.
func (es *tag) MatchOffer(offer string, lastNotificationSentDate time.Time) ([]string, error) {
	q := newTagQuery(offer, lastNotificationSentDate, "wantAddedAt")
	res, err := es.c.Search().
		Index(es.index).
		Query(q).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	matchTags := []string{}
	for _, hit := range res.Hits.Hits {
		var record types.TagESRecord
		err := json.Unmarshal(hit.Source, &record)
		if err != nil {
			return nil, err
		}
		matchTags = append(matchTags, record.Name)
	}

	return matchTags, nil
}

// MatchWant matches offers for the given want.
func (es *tag) MatchWant(want string, lastNotificationSentDate time.Time) ([]string, error) {
	q := newTagQuery(want, lastNotificationSentDate, "offerAddedAt")
	res, err := es.c.Search().
		Index(es.index).
		Query(q).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	matchTags := []string{}
	for _, hit := range res.Hits.Hits {
		var record types.TagESRecord
		err := json.Unmarshal(hit.Source, &record)
		if err != nil {
			return nil, err
		}
		matchTags = append(matchTags, record.Name)
	}

	return matchTags, nil
}
