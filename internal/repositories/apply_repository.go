package repositories

import (
	"context"
	"errors"

	"github.com/Abhishekdx300/jobster/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApplyRepository struct {
	collection *mongo.Collection
}

func NewApplyRepository(db *mongo.Database) *ApplyRepository {
	return &ApplyRepository{
		collection: db.Collection("applies"),
	}
}

func (r *ApplyRepository) FindByJobId(ctx context.Context, jobId string) ([]models.Apply, error) {
	jobObjId, err := primitive.ObjectIDFromHex(jobId)
	if err != nil {
		return nil, errors.New("invalid job id format")
	}

	cursor, err := r.collection.Find(ctx, bson.M{"job_id": jobObjId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var applies []models.Apply
	if err = cursor.All(ctx, &applies); err != nil {
		return nil, err
	}

	return applies, err
}

func (r *ApplyRepository) Create(ctx context.Context, jobId string, apply *models.Apply) (*models.Apply, error) {
	apply.ID = primitive.NewObjectID()
	jobObjId, err := primitive.ObjectIDFromHex(jobId)
	if err != nil {
		return nil, errors.New("invalid job id format")
	}

	apply.JobId = jobObjId

	_, err = r.collection.InsertOne(ctx, apply)
	if err != nil {
		return nil, err
	}
	return apply, nil
}

func (r *ApplyRepository) Update(ctx context.Context, id string, apply models.Apply) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid apply id format")
	}

	update := bson.M{
		"$set": bson.M{
			"role":         apply.Role,
			"date_applied": apply.DateApplied,
			"status":       apply.Status,
			"comment":      apply.Comment,
			"link":         apply.Link,
			"base_salary":  apply.BaseSalary,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("apply not found to update")
	}

	return nil
}

func (r *ApplyRepository) Delete(ctx context.Context, id string) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objId})
	return err
}

func (r *ApplyRepository) DeleteByJobID(ctx context.Context, jobID string) error {
	jobObjID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return errors.New("invalid job id format")
	}

	_, err = r.collection.DeleteMany(ctx, bson.M{"job_id": jobObjID})
	return err
}

func (r *ApplyRepository) CountByJobID(ctx context.Context, jobID string) (int64, error) {
	jobObjID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return 0, errors.New("invalid job id format")
	}
	count, err := r.collection.CountDocuments(ctx, bson.M{"job_id": jobObjID})
	if err != nil {
		return 0, err
	}

	return count, nil
}
