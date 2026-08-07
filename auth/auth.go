package auth

import (
	"crypto/rsa"
	"errors"
	"log"
	"net/http"
	"strings"
	"crypto/sha256"
	"encoding/hex"

	"github.com/USACE/pallid_sturgeon_api/server/models"
	"github.com/USACE/pallid_sturgeon_api/server/stores"
	"github.com/golang-jwt/jwt/v4"
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

func (a *Auth) AuthorizeAdminOrSelf(handler echo.HandlerFunc) echo.HandlerFunc {
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
		if role.Role == "ADMINISTRATOR" || user.Email == c.Param("email") {
			return handler(c)
		}

		return echo.NewHTTPError(http.StatusUnauthorized, "")
	}
}

func (a *Auth) AuthorizeViaToken(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		access := c.Request().Header.Get("access-key")
		secret := c.Request().Header.Get("secret-key")
		if access == "" {
			access = c.QueryParam("access")
			secret = c.QueryParam("secret")
		}

		if len(access) == 0 {
			return echo.NewHTTPError(http.StatusUnauthorized, "Access/Secret must be provided via query param access and secret or headers access-key and secret-key")
		}

		tokenInfo, err := a.Store.GetTokenByAccess(access)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Token Fetch via Access Key Failed", err))
		}

		if tokenInfo.TokenExpiration != "Valid" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token Expired. Please access the Pallid Sturgeon website and refresh your token.")
		}

		if HashValid(secret, tokenInfo.TokenSecret) {
			return handler(c)
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

func HashValid(input string, hash string) bool {
	
	// Convert text to bytes and compute the raw 32-byte array
	bytes := []byte(input)
	hashArray := sha256.Sum256(bytes) 

	// Convert the 32-byte array into a readable 64-character hex string
	hexString := hex.EncodeToString(hashArray[:])
	// fmt.Println("input: "+input)
	// fmt.Println("hash: "+hash)
	// fmt.Println("hexString: " +hexString)

	return hexString == hash
}