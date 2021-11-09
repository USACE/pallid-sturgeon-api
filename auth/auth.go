package auth

import (
	"crypto/rsa"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/USACE/pallid_sturgeon_api/server/models"
	"github.com/USACE/pallid_sturgeon_api/server/stores"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

const (
	PUBLIC = iota
	ADMIN
	OFFICEADMIN
	OFFICEUSER
	READONLY
)

type Auth struct {
	Store     *stores.AuthStore
	VerifyKey *rsa.PublicKey
}

/*
Authorize Options:
1) Public - All KEYCLOAK Users
2) PM - Project Manager
3) TM - Team Member
*/

func (a *Auth) Authorize(handler echo.HandlerFunc, roles ...int) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get(echo.HeaderAuthorization)
		tokenString := strings.TrimPrefix(auth, "Bearer ")
		claims, err := a.marshalJwt(tokenString)
		//claims, err := marshalJwts(tokenString)
		if err != nil {
			log.Print(err)
			return echo.NewHTTPError(http.StatusUnauthorized, "bad token")
		}
		user, err := a.Store.GetUserFromJwt(claims)
		if err != nil {
			return err
		}
		role, err := a.Store.GetUserRoleOffice(user.Email)
		if err != nil {
			return err
		}
		c.Set("PSUSER", user)
		switch {
		case contains(roles, PUBLIC):
			return handler(c)
		case contains(roles, ADMIN):
			if role.Role == "ADMINISTRATOR" {
				return handler(c)
			}
		case contains(roles, OFFICEADMIN):
			if role.Role == "OFFICE ADMIN" {
				return handler(c)
			}
		case contains(roles, OFFICEUSER):
			if role.Role == "OFFICE USER" {
				return handler(c)
			}
		case contains(roles, READONLY):
			if role.Role == "READONLY" {
				return handler(c)
			}
		}
		return echo.NewHTTPError(http.StatusUnauthorized, "")
	}
}

func (a *Auth) LoadVerificationKey(publicKey string) error {
	pk, err := jwt.ParseRSAPublicKeyFromPEM([]byte("-----BEGIN PUBLIC KEY-----\n" + publicKey + "\n-----END PUBLIC KEY-----"))
	if err != nil {
		return err
	}
	a.VerifyKey = pk
	return nil
}

func (a *Auth) marshalJwt(tokenString string) (models.JwtClaim, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return a.VerifyKey, nil
	})
	if err != nil {
		return models.JwtClaim{}, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		jwtUser := models.JwtClaim{
			//CacUid:    claims["cacUID"].(string),
			Name:      claims["name"].(string),
			Email:     claims["email"].(string),
			FirstName: claims["given_name"].(string),
			LastName:  claims["family_name"].(string),
		}
		return jwtUser, nil
	} else {
		return models.JwtClaim{}, errors.New("Invalid Token")
	}
}

func contains(a []int, x int) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
