package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name           string             `bson:"name" json:"name"`
	DateAdded      time.Time          `bson:"date_added" json:"date_added"`
	PersonalRating int                `bson:"personal_rating,omitempty" json:"personal_rating,omitempty"`
	Tags           []string           `bson:"tags,omitempty" json:"tags,omitempty"`
}
