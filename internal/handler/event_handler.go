package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"user-activity-tracking-service/internal/models"
	"user-activity-tracking-service/internal/service"
)

const MaxRequestBodyBytes = 64 * 1024 // 64 KB

// EventHandler handles HTTP requests for events.
type EventHandler struct {
	service service.EventService
}

// NewEventHandler creates a new EventHandler instance.
func NewEventHandler(service service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

// RegisterRoutes binds event endpoints to the provided HTTP ServeMux.
func (h *EventHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/events", h.IngestEvent)
	mux.HandleFunc("GET /api/v1/events", h.ListEvents)
}

// IngestEvent handles POST /api/v1/events.
func (h *EventHandler) IngestEvent(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodyBytes)

	var req models.IngestEventRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			respondError(w, http.StatusRequestEntityTooLarge, "request body exceeds maximum limit of 64 KB")
			return
		}
		respondError(w, http.StatusBadRequest, fmt.Sprintf("invalid JSON payload: %v", err))
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	event, err := h.service.IngestEvent(r.Context(), req)
	if err != nil {
		if errors.Is(err, models.ErrInvalidUserID) ||
			errors.Is(err, models.ErrEmptyAction) ||
			errors.Is(err, models.ErrActionTooLong) ||
			errors.Is(err, models.ErrInvalidMetadata) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondError(w, http.StatusInternalServerError, "failed to ingest event")
		return
	}

	respondJSON(w, http.StatusCreated, event)
}

// ListEvents handles GET /api/v1/events.
func (h *EventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := models.ListEventsFilter{
		Limit:  50,
		Offset: 0,
	}

	if userIDStr := query.Get("user_id"); userIDStr != "" {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil || userID <= 0 {
			respondError(w, http.StatusBadRequest, "invalid 'user_id' parameter: must be a positive integer")
			return
		}
		filter.UserID = &userID
	}

	if fromStr := query.Get("from"); fromStr != "" {
		fromTime, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid 'from' parameter: must be RFC3339 timestamp (e.g. 2026-08-01T00:00:00Z)")
			return
		}
		filter.From = &fromTime
	}

	if toStr := query.Get("to"); toStr != "" {
		toTime, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid 'to' parameter: must be RFC3339 timestamp (e.g. 2026-08-15T23:59:59Z)")
			return
		}
		filter.To = &toTime
	}

	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			respondError(w, http.StatusBadRequest, "invalid 'limit' parameter: must be a positive integer")
			return
		}
		filter.Limit = limit
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			respondError(w, http.StatusBadRequest, "invalid 'offset' parameter: must be a non-negative integer")
			return
		}
		filter.Offset = offset
	}

	events, err := h.service.GetEvents(r.Context(), filter)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDateRange) || errors.Is(err, models.ErrInvalidUserID) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondError(w, http.StatusInternalServerError, "failed to retrieve events")
		return
	}

	respondJSON(w, http.StatusOK, events)
}
