package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Abhishekdx300/jobster/internal/constants"
	"github.com/Abhishekdx300/jobster/internal/helpers"
	"github.com/Abhishekdx300/jobster/internal/models"
	"github.com/Abhishekdx300/jobster/internal/services"
	"github.com/go-chi/chi/v5"
)

type ApplyHandler struct {
	service *services.ApplyService
}

func NewApplyHandler(s *services.ApplyService) *ApplyHandler {
	return &ApplyHandler{service: s}
}

func (h *ApplyHandler) FindByJobId(w http.ResponseWriter, r *http.Request) {
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

func (h *ApplyHandler) Create(w http.ResponseWriter, r *http.Request) {
	jobId := chi.URLParam(r, "jobId")

	var apply models.Apply
	if err := json.NewDecoder(r.Body).Decode(&apply); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate
	if !constants.IsValidApplyTag(apply.Status) {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid status option"))
		return
	}

	// role, comment, link
	isValid := helpers.ValidateSize(200, len(apply.Role), len(apply.Comment), len(apply.Link))
	if !isValid {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("Input size too long"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	createdApply, err := h.service.Create(ctx, jobId, &apply)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, createdApply)
}

func (h *ApplyHandler) Update(w http.ResponseWriter, r *http.Request) {
	jobId := chi.URLParam(r, "jobId")
	id := chi.URLParam(r, "applyId")
	var apply models.Apply
	if err := json.NewDecoder(r.Body).Decode(&apply); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate
	if !constants.IsValidApplyTag(apply.Status) {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid status option"))
		return
	}

	// role, comment, link

	// role, comment, link
	isValid := helpers.ValidateSize(200, len(apply.Role), len(apply.Comment), len(apply.Link))
	if !isValid {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("Input size too long"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.service.Update(ctx, jobId, id, apply); err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Apply updated successfully"})
}

func (h *ApplyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	jobId := chi.URLParam(r, "jobId")
	id := chi.URLParam(r, "applyId")

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.service.Delete(ctx, jobId, id); err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Apply deleted successfully"})
}
