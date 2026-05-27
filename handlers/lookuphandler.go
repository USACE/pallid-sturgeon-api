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
	years, err := s.Store.GetYears()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve years data", err))
	}

	fieldOffices, err := s.Store.GetFieldOffices()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve field offices data", err))
	}

	fieldOfficeSegments, err := s.Store.GetFieldOfficeSegment()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve field office segments data", err))
	}

	projects, err := s.Store.GetProjects()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve projects data", err))
	}

	segments, err := s.Store.GetSegments()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve segments data", err))
	}

	seasons, err := s.Store.GetSeasons()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve seasons data", err))
	}

	sampleUnitTypes, err := s.Store.GetSampleUnitTypes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve sample unit types data", err))
	}

	bendSelections, err := s.Store.GetBendSelections()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve bend selections data", err))
	}

	bendRiverMile, err := s.Store.GetBendRiverMile()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve bend river mile data", err))
	}

	chutes, err := s.Store.GetChutes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve chutes data", err))
	}

	reach, err := s.Store.GetReach()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve reach data", err))
	}

	// MORIVER LOOK UP QUERIES
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

	subsampleTypes, err := s.Store.GetSubsampleTypes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve subsample types data", err))
	}

	// FISH LOOK UP QUERIES
	fishCodes, err := s.Store.GetFishCodes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve fish codes data", err))
	}

	fishStructures, err := s.Store.GetFishStructures()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve fish structures data", err))
	}

	floyTagPrefixes, err := s.Store.GetFloyTagPrefixes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve floy tag prefixes data", err))
	}

	lengthTypes, err := s.Store.GetLengthTypes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve length types data", err))
	}

	markRecaptureOptions, err := s.Store.GetMarkRecapture()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve mark recapture data", err))
	}

	frequencyId, err := s.Store.GetFrequencyId()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve frequency id data", err))
	}

	// SUPPLEMENTAL LOOK UP MODELS
	pitRnzOptions, err := s.Store.GetPitRnzOptions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve pit RNZ option data", err))
	}

	elastomerColorOptions, err := s.Store.GetElastomerColors()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve elastomer color data", err))
	}

	elastomerHvxOptions, err := s.Store.GetElastomerHvxOptions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve elastomer HVX option data", err))
	}

	pallidLocationStatusOptions, err := s.Store.GetPallidLocationStatusOptions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve pallid location status option data", err))
	}

	hatcheryOriginOptions, err := s.Store.GetHatcheryOriginOptions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve hatchery origin data", err))
	}
	
	// PROCEDURE LOOK UP MODELS
	
	purposeOptions, err := s.Store.GetPurposeOptions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve purpose option data", err))
	}
	
	evalLocationOptions, err := s.Store.GetEvalLocationOptions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve evaluation location option data", err))
	}

	sexOptions, err := s.Store.GetSexOptions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve sex option data", err))
	}

	reproductiveStatusOptions, err := s.Store.GetReproductiveStatusOptions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve reproductive status option data", err))
	}

	spawnBehavior, err := s.Store.GetSpawnBehavior()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve spawn behavior data", err))
	}

	positionConfidence, err := s.Store.GetPositionConfidence()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve position confidence data", err))
	}

	searchTypeCodes, err := s.Store.GetSearchType()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve search type data", err))
	}
	// GENERAL LOOK UP QUERIES
	yesNoOptions, err := s.Store.GetYesNoOptions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve yes/no option data", err))
	}


	frequencyIds, err := s.Store.GetFrequencyIds()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to retrieve Frequency ID data", err))
	}

	// Single combined response
	response := map[string]any{
		"bendSelections":    bendSelections,
		"bendRiverMile": 	 bendRiverMile,
		"chutes": 			 chutes,
		"estimations": 		 estimations,
		"fieldOffices": 	 fieldOffices,
		"fieldOfficeSegments": fieldOfficeSegments,
		"filteredGearCodes": filteredGearCodes,
		"fishCodes":         fishCodes,
		"fishStructures":    fishStructures,
		"floyTagPrefixes":   floyTagPrefixes,
		"gearCodes":         gearCodes,
		"gearTypes":         gearTypes,
		"lengthTypes": 	     lengthTypes,
		"macros":            macros,
		"mesos":             mesos,
		"macroMesos":        macroMesos,
		"markRecaptureOptions": markRecaptureOptions,
		"microSetSite": 	 microSetSite,
		"microStructures": 	 microStructures,
		"microHabitats":     microHabitats,
		"projects": 		 projects,
		"reach": 			 reach,
		"sampleUnitTypes": 	 sampleUnitTypes,
		"seasons": 			 seasons,
		"segments": 		 segments,
		"setSite1Options":   setSite1,
		"setSite2Options":   setSite2,
		"setSite3Options":   setSite3,
		"structureFlows": 	 structureFlows,
		"structureMods": 	 structureMods,
		"subsampleTypes":    subsampleTypes,
		"u6Options":         u6,
		"u7Options":         u7,
		"years": 			 years,
		"frequencyId":		 frequencyId,
		"spawnBehavior":	 spawnBehavior,
		"positionConfidence": positionConfidence,
		"searchTypeCodes":	 searchTypeCodes,
		"frequencyIds":    frequencyIds,
		"pitRnzOptions": pitRnzOptions,
		"elastomerColorOptions": elastomerColorOptions,
		"elastomerHvxOptions": elastomerHvxOptions,
		"pallidLocationStatusOptions": pallidLocationStatusOptions,
		"hatcheryOriginOptions": hatcheryOriginOptions,
		"purposeOptions": purposeOptions,
		"evalLocationOptions": evalLocationOptions,
		"sexOptions": sexOptions,
		"reproductiveStatusOptions": reproductiveStatusOptions,
		"yesNoOptions": yesNoOptions,
	}

	return c.JSON(http.StatusOK, models.NewSuccessResponse("Lookup data retrieved successfully", response))
}
