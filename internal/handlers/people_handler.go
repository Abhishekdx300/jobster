package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Abhishekdx300/jobster/internal/helpers"
	"github.com/Abhishekdx300/jobster/internal/models"
	"github.com/Abhishekdx300/jobster/internal/services"
	"github.com/go-chi/chi/v5"
)

type PeopleHandler struct {
	service *services.PeopleService
}

func NewPeopleHandler(s *services.PeopleService) *PeopleHandler {
	return &PeopleHandler{service: s}
}

func (h *PeopleHandler) FindByJobId(w http.ResponseWriter, r *http.Request) {
	jobId := chi.URLParam(r, "jobId")

	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	applies, err := h.service.FindByJobId(ctx, jobId)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, applies)
}

func (h *PeopleHandler) Create(w http.ResponseWriter, r *http.Request) {
	jobId := chi.URLParam(r, "jobId")

	var peopleReach models.PeopleReach
	if err := json.NewDecoder(r.Body).Decode(&peopleReach); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	isValid := helpers.ValidateSize(50, len(peopleReach.ProfileLink), len(peopleReach.Comment))
	if !isValid {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("Input size too long"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	createdPeopleReach, err := h.service.Create(ctx, jobId, &peopleReach)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, createdPeopleReach)
}

func (h *PeopleHandler) Update(w http.ResponseWriter, r *http.Request) {
	jobId := chi.URLParam(r, "jobId")
	id := chi.URLParam(r, "peopleId")
	var peopleReach models.PeopleReach
	if err := json.NewDecoder(r.Body).Decode(&peopleReach); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	isValid := helpers.ValidateSize(50, len(peopleReach.ProfileLink), len(peopleReach.Comment))
	if !isValid {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("Input size too long"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.service.Update(ctx, jobId, id, peopleReach); err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "people reach updated successfully"})
}

func (h *PeopleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	jobId := chi.URLParam(r, "jobId")
	id := chi.URLParam(r, "peopleId")

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.service.Delete(ctx, jobId, id); err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "people reach deleted successfully"})
}
