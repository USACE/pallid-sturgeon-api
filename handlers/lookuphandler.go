package handlers

import (
	"net/http"

	"github.com/USACE/pallid_sturgeon_api/server/models"
	"github.com/USACE/pallid_sturgeon_api/server/stores"
	"github.com/labstack/echo/v4"
)

type LookupHandler struct {
	Store *stores.LookupStore
}

func (s *LookupHandler) GetAllLookups(c echo.Context) error {
	// GLOBAL LOOK UP QUERIES
	frequencyIds, err := s.Store.GetFrequencyIds()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve Frequency ID data", err))
	}

	scuteLocations, err := s.Store.GetScuteLocations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve Scute Location data", err))
	}

	// Single combined response
	response := map[string]any{
		"frequencyIds":    frequencyIds,
		"scuteLocations":  scuteLocations,
	}

	return c.JSON(http.StatusOK, models.NewSuccessResponse("Lookup data retrieved successfully", response))
}