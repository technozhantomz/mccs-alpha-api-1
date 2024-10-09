package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	DeletedAt time.Time          `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`

	Name         string    `json:"name,omitempty" bson:"name,omitempty"`
	OfferAddedAt time.Time `json:"offerAddedAt,omitempty" bson:"offerAddedAt,omitempty"`
	WantAddedAt  time.Time `json:"wantAddedAt,omitempty" bson:"wantAddedAt,omitempty"`
}

// Helper functions

func TagToNames(tags []*Tag) []string {
	names := make([]string, 0, len(tags))
	for _, t := range tags {
		names = append(names, t.Name)
	}
	return names
}

// Helper types

type FindTagResult struct {
	Tags            []*Tag
	NumberOfResults int
	TotalPages      int
}

type MatchedTags struct {
	MatchedOffers map[string][]string `json:"matchedOffers,omitempty"`
	MatchedWants  map[string][]string `json:"matchedWants,omitempty"`
}
