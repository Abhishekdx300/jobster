package services

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/Abhishekdx300/jobster/internal/constants"
	"github.com/Abhishekdx300/jobster/internal/models"
	"github.com/Abhishekdx300/jobster/internal/repositories"
)

type JobService struct {
	jobRepo    *repositories.JobRepository
	applyRepo  *repositories.ApplyRepository
	peopleRepo *repositories.PeopleRepository
}

func NewJobService(jobRepo *repositories.JobRepository, applyRepo *repositories.ApplyRepository, peopleRepo *repositories.PeopleRepository) *JobService {
	return &JobService{
		jobRepo:    jobRepo,
		applyRepo:  applyRepo,
		peopleRepo: peopleRepo,
	}
}

func (s *JobService) Search(ctx context.Context, params repositories.SearchParams) ([]models.Job, int64, error) {
	return s.jobRepo.SearchAndFilter(ctx, params)
}

func (s *JobService) GetById(ctx context.Context, id string) (*models.Job, error) {
	return s.jobRepo.FindById(ctx, id)
}

func (s *JobService) GetAll(ctx context.Context) ([]models.Job, error) {
	return s.jobRepo.GetAll(ctx)
}

func (s *JobService) Create(ctx context.Context, job *models.Job) (*models.Job, error) {
	// db validations

	// check name already available
	jobName := strings.TrimSpace(job.Name)
	alreadyPresent, _ := s.jobRepo.FindByName(ctx, jobName)

	if alreadyPresent != nil {
		return nil, errors.New("Job already present with the name.")
	}
	// else create

	for _, tag := range job.Tags {
		if !constants.IsValidJobTag(tag) {
			return nil, errors.New("invalid tag provided: " + tag)
		}
	}

	job.Name = jobName
	job.DateAdded = time.Now()

	return s.jobRepo.Create(ctx, job)
}

func (s *JobService) Update(ctx context.Context, id string, job *models.Job) error {
	for _, tag := range job.Tags {
		if !constants.IsValidJobTag(tag) {
			return errors.New("invalid tag provided: " + tag)
		}
	}
	return s.jobRepo.Update(ctx, id, job)
}

func (s *JobService) Delete(ctx context.Context, id string) error {

	err := s.applyRepo.DeleteByJobID(ctx, id)
	if err != nil {
		log.Println(err.Error())
	}

	err = s.peopleRepo.DeleteByJobID(ctx, id)
	if err != nil {
		log.Println(err.Error())
	}

	return s.jobRepo.Delete(ctx, id)
}
