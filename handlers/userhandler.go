package handlers

import (
	"net/http"

	"github.com/USACE/pallid_sturgeon_api/server/models"
	"github.com/USACE/pallid_sturgeon_api/server/stores"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Store *stores.AuthStore
}

func (u *UserHandler) AddUserRoleOffice(c echo.Context) error {
	var err error
	uUserRoleOffice := models.UserRoleOffice{}
	if err := c.Bind(&uUserRoleOffice); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = u.Store.AddUserRoleOffice(uUserRoleOffice)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (u *UserHandler) GetUserRoleOffice(c echo.Context) error {
	email := c.Param("email")

	roleOffice, err := u.Store.GetUserRoleOffice(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, roleOffice)
}

func (u *UserHandler) GetUsers(c echo.Context) error {
	users, err := u.Store.GetUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

func (u *UserHandler) GetUserAccessRequests(c echo.Context) error {
	users, err := u.Store.GetUserAccessRequests()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}
