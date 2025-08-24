package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Apply struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	JobId       primitive.ObjectID `bson:"job_id" json:"job_id"`
	Role        string             `bson:"role" json:"role"`
	DateApplied time.Time          `bson:"date_applied" json:"date_applied"`
	Status      string             `bson:"status" json:"status"`
	Comment     string             `bson:"comment,omitempty" json:"comment,omitempty"`
	Link        string             `bson:"link,omitempty" json:"link,omitempty"`
	BaseSalary  int                `bson:"base_salary,omitempty" json:"base_salary,omitempty"`
}
