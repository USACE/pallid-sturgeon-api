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
	bendSelections, err := s.Store.GetBendSelections()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve bend selections", err))
	}

	gearCodes, err := s.Store.GetGearCodes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve gear codes", err))
	}

	filteredGearCodes, err := s.Store.GetFilteredGearCodes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve filtered gear code data", err))
	}

	gearTypes, err := s.Store.GetGearTypes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve gear types", err))
	}

	macros, err := s.Store.GetMacros()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve macros", err))
	}

	mesos, err := s.Store.GetMesos()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve mesos", err))
	}

	macroMesos, err := s.Store.GetMacroMesos()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve macros/mesos data", err))
	}

	microStructures, err := s.Store.GetMicroStructures()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve micro structures", err))
	}

	structureFlows, err := s.Store.GetStructureFlow()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve structure flow data", err))
	}

	structureMods, err := s.Store.GetStructureMod()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve structure mod data", err))
	}

	microHabitats, err := s.Store.GetMicroHabitats()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve micro habitats", err))
	}

	u6, err := s.Store.GetU6()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve U6 data", err))
	}

	u7, err := s.Store.GetU7()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve U7 data", err))
	}

	estimations, err := s.Store.GetEstimations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve estimation data", err))
	}

	microSetSite, err := s.Store.GetMicroSetSite()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve micro & set site data", err))
	}

	setSite1, err := s.Store.GetSetSite1()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve set site 1 data", err))
	}

	setSite2, err := s.Store.GetSetSite2()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve set site 2 data", err))
	}

	setSite3, err := s.Store.GetSetSite3()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve set site 3 data", err))
	}

	bendRiverMile, err := s.Store.GetBendRiverMile()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve bend river mile data", err))
	}

	// Single combined response
	response := map[string]any{
		"bendSelections":    bendSelections,
		"bendRiverMile": 	 bendRiverMile,
		"estimations": 		 estimations,
		"filteredGearCodes": filteredGearCodes,
		"gearCodes":         gearCodes,
		"gearTypes":         gearTypes,
		"macros":            macros,
		"mesos":             mesos,
		"macroMesos":        macroMesos,
		"microSetSite": 	 microSetSite,
		"microStructures": 	 microStructures,
		"microHabitats":     microHabitats,
		"setSite1Options":   setSite1,
		"setSite2Options":   setSite2,
		"setSite3Options":   setSite3,
		"structureFlows": 	 structureFlows,
		"structureMods": 	 structureMods,
		"u6Options":         u6,
		"u7Options":         u7,
	}

	return c.JSON(http.StatusOK, models.NewSuccessResponse("Lookup data retrieved successfully", response))
}
