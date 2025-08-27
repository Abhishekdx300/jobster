package services

import (
	"context"
	"errors"

	"github.com/Abhishekdx300/jobster/internal/models"
	"github.com/Abhishekdx300/jobster/internal/repositories"
)

const applyLimit = 5

type ApplyService struct {
	applyRepo *repositories.ApplyRepository
	jobRepo   *repositories.JobRepository
}

func NewApplyService(applyRepo *repositories.ApplyRepository, jobRepo *repositories.JobRepository) *ApplyService {
	return &ApplyService{applyRepo: applyRepo, jobRepo: jobRepo}
}

func (s *ApplyService) FindByJobId(ctx context.Context, jobId string) ([]models.Apply, error) {
	// check if jobId is present or not
	_, err := s.jobRepo.FindById(ctx, jobId)
	if err != nil {
		return nil, errors.New("job not found")
	}

	return s.applyRepo.FindByJobId(ctx, jobId)
}

func (s *ApplyService) Create(ctx context.Context, jobId string, apply *models.Apply) (*models.Apply, error) {

	// check if jobId is present or not
	_, err := s.jobRepo.FindById(ctx, jobId)
	if err != nil {
		return nil, errors.New("job not found")
	}

	currentCount, err := s.applyRepo.CountByJobID(ctx, jobId)
	if err != nil {
		return nil, errors.New("could not verify apply count")
	}

	if currentCount >= applyLimit {
		return nil, errors.New("apply count limit exceeded")
	}

	return s.applyRepo.Create(ctx, jobId, apply)

}

func (s *ApplyService) Update(ctx context.Context, jobId string, id string, apply models.Apply) error {

	// check if jobId is present or not
	_, err := s.jobRepo.FindById(ctx, jobId)
	if err != nil {
		return errors.New("job not found")
	}

	return s.applyRepo.Update(ctx, id, apply)
}

func (s *ApplyService) Delete(ctx context.Context, jobId, id string) error {
	// check if jobId is present or not
	_, err := s.jobRepo.FindById(ctx, jobId)
	if err != nil {
		return errors.New("job not found")
	}
	return s.applyRepo.Delete(ctx, id)
}
