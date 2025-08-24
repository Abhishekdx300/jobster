package services

import (
	"context"
	"errors"

	"github.com/Abhishekdx300/jobster/internal/models"
	"github.com/Abhishekdx300/jobster/internal/repositories"
)

const peopleLimit = 10

type PeopleService struct {
	peopleRepo *repositories.PeopleRepository
	jobRepo    *repositories.JobRepository
}

func NewPeopleService(peopleRepo *repositories.PeopleRepository, jobRepo *repositories.JobRepository) *PeopleService {
	return &PeopleService{peopleRepo: peopleRepo, jobRepo: jobRepo}
}

func (s *PeopleService) FindByJobId(ctx context.Context, jobId string) ([]models.PeopleReach, error) {
	// check if jobId is present or not
	_, err := s.jobRepo.FindById(ctx, jobId)
	if err != nil {
		return nil, errors.New("job not found")
	}

	return s.peopleRepo.FindByJobId(ctx, jobId)
}

func (s *PeopleService) Create(ctx context.Context, jobId string, people *models.PeopleReach) (*models.PeopleReach, error) {

	// check if jobId is present or not
	_, err := s.jobRepo.FindById(ctx, jobId)
	if err != nil {
		return nil, errors.New("job not found")
	}

	currentCount, err := s.peopleRepo.CountByJobID(ctx, jobId)
	if err != nil {
		return nil, errors.New("could not verify apply count")
	}

	if currentCount >= peopleLimit {
		return nil, errors.New("apply count limit exceeded")
	}

	return s.peopleRepo.Create(ctx, jobId, people)

}

func (s *PeopleService) Update(ctx context.Context, jobId string, id string, people models.PeopleReach) error {
	// check if jobId is present or not
	_, err := s.jobRepo.FindById(ctx, jobId)
	if err != nil {
		return errors.New("job not found")
	}
	return s.peopleRepo.Update(ctx, id, people)
}

func (s *PeopleService) Delete(ctx context.Context, jobId, id string) error {
	// check if jobId is present or not
	_, err := s.jobRepo.FindById(ctx, jobId)
	if err != nil {
		return errors.New("job not found")
	}
	return s.peopleRepo.Delete(ctx, id)
}
