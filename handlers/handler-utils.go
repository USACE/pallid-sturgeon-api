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

func processPallidTime(st string,ty string) string {
	t := ""

	if len(st) > 0 {
		if ty == "db" && strings.Index(st, "/") == 1 {
			st = "0" + st
		}

		f := "01/02/2006"

		if ty == "app" {
			f = "2006-01-02T15:04:05Z"
		} 

		test, err := time.Parse(f, st)
		if err != nil {
			t = ""
		} else {
			t = test.Format("02-Jan-2006")
		}
	}
	return t

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

	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			return models.SearchParams{}, err
		}
	}

	orderbyString := c.QueryParam("orderby")

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
