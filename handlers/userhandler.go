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
	userRoleOffice := models.UserRoleOffice{}
	if err := c.Bind(&userRoleOffice); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = u.Store.AddUserRoleOffice(userRoleOffice)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (u *UserHandler) GetUserRoleOffices(c echo.Context) error {
	email := c.Param("email")

	roleOfficeItems, err := u.Store.GetUserRoleOffices(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, roleOfficeItems)
}

func (u *UserHandler) GetUserRoleOfficeById(c echo.Context) error {
	id := c.Param("id")

	roleOffice, err := u.Store.GetUserRoleOfficeById(id)
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

func (u *UserHandler) UpdateUserRoleOffice(c echo.Context) error {
	userData := models.UserRoleOffice{}
	if err := c.Bind(&userData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err := u.Store.UpdateUserRoleOffice(userData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (u *UserHandler) GetUsers2(c echo.Context) error {
	users, err := u.Store.GetUsers2()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}
