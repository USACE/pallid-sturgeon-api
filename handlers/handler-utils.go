package handlers

import (
	// "strconv"
	// "strings"

	// "di2e.net/cwbi/pallid_sturgeon_api/server/models"
	// "github.com/labstack/echo/v4"
)

// func marshalQuery(c echo.Context) (models.SearchParams, error) {
// 	var page int = 0
// 	var size int = 100
// 	var orderby string = "id"
// 	var filter string = ""
// 	var phaseType string = ""
// 	var phaseStatus string = ""
// 	var err error

// 	pageStr := c.QueryParam("page")

// 	if pageStr != "" {
// 		page, err = strconv.Atoi(pageStr)
// 		if err != nil {
// 			return models.SearchParams{}, err
// 		}
// 	}

// 	sizeStr := c.QueryParam("size")

// 	if sizeStr != "" {
// 		size, err = strconv.Atoi(sizeStr)
// 		if err != nil {
// 			return models.SearchParams{}, err
// 		}
// 	}

// 	ordebyString := c.QueryParam("orderby")

// 	if ordebyString != "" && ordebyString != "undefined" {
// 		orderby = ordebyString
// 	}

// 	filterStr := c.QueryParam("filter")

// 	if filterStr != "" {
// 		filter = filterStr
// 		if strings.Contains(filterStr, "project_type") {
// 			filter = strings.ReplaceAll(filter, "project_type", "projt.type_name")
// 		}
// 		if strings.Contains(filterStr, "id ILIKE") {
// 			filter = strings.ReplaceAll(filter, "id ILIKE", "CAST(p.id AS TEXT) ILIKE")
// 		}
// 		if strings.Contains(filterStr, "id IN") {
// 			filter = strings.ReplaceAll(filter, "id IN", "p.id IN")
// 		}
// 		if strings.Contains(filterStr, "fiscal_year ILIKE") {
// 			filter = strings.ReplaceAll(filter, "fiscal_year ILIKE", "CAST(fiscal_year AS TEXT) ILIKE")
// 		}
// 		if strings.Contains(filterStr, "project_type") {
// 			filter = strings.ReplaceAll(filter, "project_type", "projt.type_name")
// 		}
// 		if strings.Contains(filterStr, "phase ILIKE") {
// 			filter = strings.ReplaceAll(filter, "phase ILIKE", "pht.phase_name ILIKE")
// 		}
// 		if strings.Contains(filterStr, "phase IN") {
// 			filter = strings.ReplaceAll(filter, "phase IN", "pht.phase_name IN")
// 		}
// 		if strings.Contains(filterStr, "project_manager") {
// 			filter = strings.ReplaceAll(filter, "project_manager", "pl.user_name")
// 		}
// 		if strings.Contains(filterStr, "lead_engineer") {
// 			filter = strings.ReplaceAll(filter, "lead_engineer", "le.user_name")
// 		}
// 		if strings.Contains(filterStr, "current_task") {
// 			filter = strings.ReplaceAll(filter, "current_task", "b.task_name")
// 		}
// 		if strings.Contains(filterStr, "percent_complete ILIKE") {
// 			filter = strings.ReplaceAll(filter, "percent_complete ILIKE", "CAST(percent_complete AS TEXT) ILIKE")
// 		}
// 		if strings.Contains(filterStr, "slippage ILIKE") {
// 			filter = strings.ReplaceAll(filter, "slippage ILIKE", "CAST((b.actual_date - b.projected_end_date) AS TEXT) ILIKE")
// 		}
// 		if strings.Contains(filterStr, "slippage IN") {
// 			filter = strings.ReplaceAll(filter, "slippage IN", "(b.actual_date - b.projected_end_date) IN")
// 		}
// 	}

// 	phaseType = c.QueryParam("phaseType")
// 	phaseStatus = c.QueryParam("phaseStatus")

// 	return models.SearchParams{
// 		Page:        page,
// 		PageSize:    size,
// 		OrderBy:     orderby,
// 		Filter:      filter,
// 		PhaseType:   phaseType,
// 		PhaseStatus: phaseStatus,
// 	}, nil
// }
