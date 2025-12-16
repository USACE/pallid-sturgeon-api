package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/USACE/pallid_sturgeon_api/server/models"
	"github.com/labstack/echo/v4"
)

func processTimeString(st string) time.Time {
	t := time.Time{}
	if len(st) > 0 {
		if !strings.HasPrefix(st, "1") {
			st = "0" + st
		}
		test, err := time.Parse("01/02/2006", st)
		if err != nil {
			t = time.Time{}
		}
		t = test
	}

	return t
}

func processStringTime(st string, ty string) *string {
	t := new(string)

	if len(st) > 0 {
		f := "1/2/2006"

		if ty == "app" {
			f = "2006-01-02"
		}

		test, err := time.Parse(f, st)

		if err != nil {
			*t = ""
		} else {
			*t = test.Format("02-Jan-2006")
		}
	}
	return t
}

func DerefString(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

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
	if sizeStr == "" {
		sizeStr = c.QueryParam("pageSize")
	}
	if sizeStr == "" {
		sizeStr = c.QueryParam("page_size")
	}
	if sizeStr == "" {
		sizeStr = c.QueryParam("pagesize")
	}

	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			return models.SearchParams{}, err
		}
	}

	orderbyString := c.QueryParam("orderby")
	if orderbyString == "" {
		orderbyString = c.QueryParam("orderBy")
	}
	if orderbyString == "" {
		orderbyString = c.QueryParam("order_by")
	}

	if orderbyString != "" && orderbyString != "undefined" {
		orderby = orderbyString
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
