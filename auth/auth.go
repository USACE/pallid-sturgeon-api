package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"di2e.net/cwbi/pallid_sturgeon_api/server/models"
	"di2e.net/cwbi/pallid_sturgeon_api/server/stores"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

const (
	PUBLIC = iota
	PM
	TM
)

type Auth struct {
	Store     *stores.AuthStore
	VerifyKey *rsa.PublicKey
}

var verifyKeys []*rsa.PublicKey

/*
Authorize Options:
1) Public - All CAC Users
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
		c.Set("SDUSER", user)
		switch {
		case contains(roles, TM):
			for _, role := range claims.Roles {
				if role == "PM" || role == "TM" {
					return handler(c)
				}
			}
		case contains(roles, PM):
			for _, role := range claims.Roles {
				if role == "PM" {
					return handler(c)
				}
			}
		}
		return echo.NewHTTPError(http.StatusUnauthorized, "")
	}
}

func loadKeyFile(filePath string) (*rsa.PublicKey, error) {
	publicKeyBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
}

func (a *Auth) LoadVerificationKey(filePath string) error {
	publicKeyBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	pk, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
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
			Sub:   claims["sub"].(string),
			Name:  claims["name"].(string),
			Email: claims["email"].(string),
			Roles: claims["roles"].([]interface{}),
		}
		return jwtUser, nil
	} else {
		return models.JwtClaim{}, errors.New("Invalid Token")
	}
}

func LoadVerificationKeys(fieldPath string) error {
	files, err := ioutil.ReadDir(fieldPath)
	if err != nil {
		return err
	}
	for _, v := range files {
		if ext := filepath.Ext(v.Name()); ext == ".pem" {
			fmt.Printf("Loading Public Key: %s\n", v.Name())
			pk, err := loadKeyFile(fieldPath + "/" + v.Name())
			if err != nil {
				return err
			}
			verifyKeys = append(verifyKeys, pk)
		}
	}
	return nil
}

func marshalJwts(tokenString string) (models.JwtClaim, error) {
	var token *jwt.Token = nil
	var err error
	for _, verificationKey := range verifyKeys {
		token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return verificationKey, nil
		})
		if err == nil {
			break
		}
	}

	if token == nil {
		return models.JwtClaim{}, errors.New("Invalid Token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		jwtUser := models.JwtClaim{
			Sub:   claims["sub"].(string),
			Name:  claims["name"].(string),
			Email: claims["email"].(string),
			Roles: claims["roles"].([]interface{}),
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
