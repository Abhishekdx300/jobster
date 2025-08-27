package repositories

import (
	"context"
	"errors"

	"github.com/Abhishekdx300/jobster/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type JobRepository struct {
	collection *mongo.Collection
}

func NewJobRepository(db *mongo.Database) *JobRepository {
	return &JobRepository{
		collection: db.Collection("jobs"),
	}
}

type SearchParams struct {
	Name     string
	Tags     []string
	Rating   int64
	Page     int64
	PageSize int64
}

func (r *JobRepository) SearchAndFilter(ctx context.Context, params SearchParams) ([]models.Job, int64, error) {
	filter := bson.M{}

	if params.Name != "" {
		filter["name"] = bson.M{"$regex": primitive.Regex{Pattern: params.Name, Options: "i"}}
	}

	// ALL must match
	if len(params.Tags) > 0 {
		filter["tags"] = bson.M{"$all": params.Tags}
	}

	filter["personal_rating"] = bson.M{"$gte": params.Rating}

	skip := (params.Page - 1) * params.PageSize
	limit := params.PageSize

	findOption := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetProjection(bson.M{"applies": 0, "people_reached": 0})

	cursor, err := r.collection.Find(ctx, filter, findOption)

	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var jobs []models.Job

	if err = cursor.All(ctx, &jobs); err != nil {
		return nil, 0, err
	}

	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return jobs, totalCount, nil

}

func (repo *JobRepository) FindByName(ctx context.Context, name string) (*models.Job, error) {
	var job models.Job
	err := repo.collection.FindOne(ctx, bson.M{"name": name}).Decode(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (repo *JobRepository) FindById(ctx context.Context, id string) (*models.Job, error) {

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}

	var job models.Job

	err = repo.collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (repo *JobRepository) Create(ctx context.Context, job *models.Job) (*models.Job, error) {
	job.ID = primitive.NewObjectID()

	_, err := repo.collection.InsertOne(ctx, job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (repo *JobRepository) Update(ctx context.Context, id string, job *models.Job) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	update := bson.M{
		"$set": bson.M{
			"name":            job.Name,
			"personal_rating": job.PersonalRating,
			"tags":            job.Tags,
		},
	}

	result, err := repo.collection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("job not found to update")
	}

	return nil
}

func (repo *JobRepository) Delete(ctx context.Context, id string) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	result, err := repo.collection.DeleteOne(ctx, bson.M{"_id": objId})

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("job not found to delete")
	}

	return nil
}
