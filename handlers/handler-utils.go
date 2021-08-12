package handlers

import (
	"strconv"

	"github.com/USACE/pallid_sturgeon_api/server/models"
	"github.com/labstack/echo/v4"
)

func marshalQuery(c echo.Context) (models.SearchParams, error) {
	var page int = 0
	var size int = 20
	var orderby string = ""
	var filter string = ""
	var err error

	pageStr := c.QueryParam("page")

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return models.SearchParams{}, err
		}
	}

	sizeStr := c.QueryParam("size")

	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			return models.SearchParams{}, err
		}
	}

	ordebyString := c.QueryParam("orderby")

	if ordebyString != "" && ordebyString != "undefined" {
		orderby = ordebyString
	}

	filterString := c.QueryParam("filter")

	if filterString != "" && filterString != "undefined" {
		filter = filterString
	}

	return models.SearchParams{
		Page:     page,
		PageSize: size,
		OrderBy:  orderby,
		Filter:   filter,
	}, nil
}
