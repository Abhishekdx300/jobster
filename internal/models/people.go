package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PeopleReach struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	JobId       primitive.ObjectID `bson:"job_id" json:"job_id"`
	ProfileLink string             `bson:"profile_link" json:"profile_link"`
	Reached     bool               `bson:"reached" json:"reached"`
	Comment     string             `bson:"comment,omitempty" json:"comment,omitempty"`
}
