package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"user-activity-tracking-service/internal/models"
	"user-activity-tracking-service/internal/service"
)

// StatHandler handles HTTP requests for aggregated statistics.
type StatHandler struct {
	service service.StatService
}

// NewStatHandler creates a new StatHandler instance.
func NewStatHandler(service service.StatService) *StatHandler {
	return &StatHandler{service: service}
}

// RegisterRoutes binds stats endpoints to the provided HTTP ServeMux.
func (h *StatHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/stats", h.ListStats)
}

// ListStats handles GET /api/v1/stats.
func (h *StatHandler) ListStats(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := models.ListStatsFilter{
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

	stats, err := h.service.GetStats(r.Context(), filter)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDateRange) || errors.Is(err, models.ErrInvalidUserID) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondError(w, http.StatusInternalServerError, "failed to retrieve activity stats")
		return
	}

	respondJSON(w, http.StatusOK, stats)
}
