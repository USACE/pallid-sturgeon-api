package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"di2e.net/cwbi/pallid_sturgeon_api/server/config"
	"di2e.net/cwbi/pallid_sturgeon_api/server/models"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Config *config.AppConfig
}

var response models.KeyCloakResponse

func (h *UserHandler) GetAdminToken() error {
	var err error
	var decodedResponse *models.KeyCloakResponse

	hc := http.Client{}
	form := url.Values{}

	// Build form to POST
	form.Add("client_id", "admin-cli")
	form.Add("username", h.Config.AdminUsername)
	form.Add("password", h.Config.AdminPassword)
	form.Add("grant_type", "password")

	// Here is where the code stops working, returns an empty body
	tokenUrl := h.Config.KeycloakUrl + "/realms/master/protocol/openid-connect/token"
	req, err := http.NewRequest("POST", tokenUrl, strings.NewReader(form.Encode()))
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}

	resp, err := hc.Do(req)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&decodedResponse)

	if err != nil {
		return err
	} else {
		// We only need the Access Token
		response.AccessToken = decodedResponse.AccessToken
		response.Expires = decodedResponse.Expires
	}

	return err
}

func (h *UserHandler) GetClientId() error {

	url := h.Config.KeycloakUrl + "/admin/realms/" + h.Config.Realm + "/clients?clientId=" + h.Config.ClientName

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", "Bearer "+response.AccessToken)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	var bodyResponse []models.KeyCloakResponse
	json.Unmarshal([]byte(body), &bodyResponse)

	if len(bodyResponse) > 0 {
		response.ClientId = bodyResponse[0].ClientId
	}

	return err
}

func (h *UserHandler) GetRoleId(roleType string) (string, error) {

	url := h.Config.KeycloakUrl + "/admin/realms/" + h.Config.Realm + "/clients/" + response.ClientId + "/roles/" + roleType

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", "Bearer "+response.AccessToken)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	var role models.KeyCloakRole
	json.Unmarshal([]byte(body), &role)

	if role.RoleId != "" {
		return role.RoleId, err
	}

	return "", err
}

func (h *UserHandler) GetUsersByRoleType(c echo.Context) error {
	roleType := c.Param("roleType")
	var err error

	err = h.GetAdminToken()
	if err != nil {
		return err
	}

	if response.ClientId == "" {
		err = h.GetClientId()
		if err != nil {
			return err
		}
	}

	url := h.Config.KeycloakUrl + "/admin/realms/" + h.Config.Realm + "/clients/" + response.ClientId + "/roles/" + roleType + "/users"

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", "Bearer "+response.AccessToken)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	var users []models.KeyCloakUser
	json.Unmarshal([]byte(body), &users)

	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserByUsername(userName string) ([]models.KeyCloakUser, error) {
	var users []models.KeyCloakUser

	url := h.Config.KeycloakUrl + "/admin/realms/" + h.Config.Realm + "/users?username=" + userName

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return users, err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", "Bearer "+response.AccessToken)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	json.Unmarshal([]byte(body), &users)

	return users, err
}

func (h *UserHandler) AddUserRoleRequest(c echo.Context) error {
	var err error

	user := models.KeyCloakUser{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	err = h.GetAdminToken()
	if err != nil {
		return err
	}

	if response.ClientId == "" {
		err = h.GetClientId()
		if err != nil {
			return err
		}
	}

	err = h.AddUserOrUserRole(user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (h *UserHandler) AddUserOrUserRole(user models.KeyCloakUser) error {
	var err error

	roleId, err := h.GetRoleId(user.Role.RoleName)
	if err != nil {
		return err
	}
	keycloakUser, err := h.GetUserByUsername(user.Username)
	if err != nil {
		return err

	}
	if len(keycloakUser) == 0 {
		log.Println("user doesn't exist, create user")
		err = h.AddUser(user)
		if err != nil {
			return err
		}

		err = h.AddUserOrUserRole(user)
		if err != nil {
			return err
		}
	} else {
		err = h.AddUserRole(keycloakUser[0], roleId, user.Role.RoleName)
		if err != nil {
			return err
		}
	}

	return err
}

func (h *UserHandler) AddUser(user models.KeyCloakUser) error {
	var err error

	hc := http.Client{}

	values := map[string]interface{}{"firstName": user.FirstName,
		"lastName": user.LastName,
		"username": user.Username,
		"email":    user.Username,
		"enabled":  "true",
		// "credentials": []interface{}{map[string]string{"type": "password", "value": "password"}},
		"attributes": map[string]string{"cacUID": user.UserID}}
	json_data, err := json.Marshal(values)
	if err != nil {
		return err
	}

	tokenUrl := h.Config.KeycloakUrl + "/admin/realms/" + h.Config.Realm + "/users/"
	req, err := http.NewRequest("POST", tokenUrl, bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+response.AccessToken)

	resp, err := hc.Do(req)
	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(resp.Body)

	return err
}

func (h *UserHandler) AddUserRole(user models.KeyCloakUser, roleId string, roleType string) error {
	var err error

	hc := http.Client{}

	values := []interface{}{map[string]string{"id": roleId, "name": roleType}}
	json_data, err := json.Marshal(values)
	if err != nil {
		return err
	}

	// Here is where the code stops working, returns an empty body
	tokenUrl := h.Config.KeycloakUrl + "/admin/realms/" + h.Config.Realm + "/users/" + user.UserID + "/role-mappings/clients/" + response.ClientId
	req, err := http.NewRequest("POST", tokenUrl, bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+response.AccessToken)

	resp, err := hc.Do(req)
	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(resp.Body)

	return err
}

func (h *UserHandler) DeleteUserRole(c echo.Context) error {
	userId := c.Param("userId")
	roleType := c.Param("roleType")
	var err error

	hc := http.Client{}

	roleId, err := h.GetRoleId(roleType)
	if err != nil {
		return err
	}

	values := []interface{}{map[string]string{"id": roleId, "name": roleType}}
	json_data, err := json.Marshal(values)
	if err != nil {
		return err
	}

	tokenUrl := h.Config.KeycloakUrl + "/admin/realms/" + h.Config.Realm + "/users/" + userId + "/role-mappings/clients/" + response.ClientId
	req, err := http.NewRequest("DELETE", tokenUrl, bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+response.AccessToken)

	resp, err := hc.Do(req)
	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(resp.Body)

	return err
}
