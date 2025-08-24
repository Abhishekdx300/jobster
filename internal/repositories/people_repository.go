package repositories

import (
	"context"
	"errors"

	"github.com/Abhishekdx300/jobster/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PeopleRepository struct {
	collection *mongo.Collection
}

func NewPeopleRepository(db *mongo.Database) *PeopleRepository {
	return &PeopleRepository{
		collection: db.Collection("people_reach"),
	}
}

func (r *PeopleRepository) FindByJobId(ctx context.Context, jobId string) ([]models.PeopleReach, error) {
	jobObjId, err := primitive.ObjectIDFromHex(jobId)
	if err != nil {
		return nil, errors.New("invalid job id format")
	}

	cursor, err := r.collection.Find(ctx, bson.M{"job_id": jobObjId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var peopleReaches []models.PeopleReach
	if err = cursor.All(ctx, &peopleReaches); err != nil {
		return nil, err
	}
	return peopleReaches, err
}

func (r *PeopleRepository) Create(ctx context.Context, jobId string, peopleReach *models.PeopleReach) (*models.PeopleReach, error) {
	peopleReach.ID = primitive.NewObjectID()
	jobObjId, err := primitive.ObjectIDFromHex(jobId)
	if err != nil {
		return nil, errors.New("invalid job id format")
	}

	peopleReach.JobId = jobObjId

	_, err = r.collection.InsertOne(ctx, peopleReach)
	if err != nil {
		return nil, err
	}
	return peopleReach, nil
}

func (r *PeopleRepository) Update(ctx context.Context, id string, peopleReach models.PeopleReach) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid people reach id format")
	}

	update := bson.M{
		"$set": bson.M{
			"profile_link": peopleReach.ProfileLink,
			"reached":      peopleReach.Reached,
			"comment":      peopleReach.Comment,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("people reach not found to update")
	}

	return nil
}

func (r *PeopleRepository) Delete(ctx context.Context, id string) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objId})
	return err
}

func (r *PeopleRepository) DeleteByJobID(ctx context.Context, jobID string) error {
	jobObjID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return errors.New("invalid job id format")
	}

	_, err = r.collection.DeleteMany(ctx, bson.M{"job_id": jobObjID})
	return err
}

func (r *PeopleRepository) CountByJobID(ctx context.Context, jobID string) (int64, error) {
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
