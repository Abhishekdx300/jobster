package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Abhishekdx300/jobster/internal/helpers"
	"github.com/Abhishekdx300/jobster/internal/models"
	"github.com/Abhishekdx300/jobster/internal/repositories"
	"github.com/Abhishekdx300/jobster/internal/services"
	"github.com/go-chi/chi/v5"
)

type JobHandler struct {
	service *services.JobService
}

func NewJobHandler(s *services.JobService) *JobHandler {
	return &JobHandler{service: s}
}

func (h *JobHandler) Search(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	name := queryParams.Get("name")
	var tags []string
	tagsQuery := queryParams.Get("tags")
	if tagsQuery != "" {
		tags = strings.Split(tagsQuery, ",")
	}

	rating, err := strconv.ParseInt(queryParams.Get("rating"), 10, 64)
	if err != nil || rating < 1 || rating > 5 {
		rating = 1
	}

	page, err := strconv.ParseInt(queryParams.Get("page"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(queryParams.Get("pageSize"), 10, 64)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	params := repositories.SearchParams{
		Name:     name,
		Tags:     tags,
		Rating:   rating,
		Page:     page,
		PageSize: pageSize,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	jobs, totalCount, err := h.service.Search(ctx, params)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]any{
		"data": jobs,
		"pagination": map[string]any{
			"total":      totalCount,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": (totalCount + pageSize - 1) / pageSize,
		},
	}

	helpers.WriteJSON(w, http.StatusOK, response)
}

func (h *JobHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	jobs, err := h.service.GetAll(ctx)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, jobs)
}

func (h *JobHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "jobId")
	// can modify context for timeout
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	job, err := h.service.GetById(ctx, id)
	if err != nil {
		helpers.WriteError(w, http.StatusNotFound, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, job)
}

func (h *JobHandler) Create(w http.ResponseWriter, r *http.Request) {
	var job models.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	createdJob, err := h.service.Create(ctx, &job)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, createdJob)

}

func (h *JobHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "jobId")
	var job models.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.service.Update(ctx, id, &job); err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Job updated successfully"})
}

func (h *JobHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "jobId")

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.service.Delete(ctx, id); err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Job deleted successfully"})
}
