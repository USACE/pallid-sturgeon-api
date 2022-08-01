package models

type JwtClaim struct {
	CacUid    *string
	Name      string
	Email     string
	FirstName string
	LastName  string
	Roles     []interface{}
}

type SearchParams struct {
	Page     int    `json:"page"`
	PageSize int    `json:"size"`
	OrderBy  string `json:"orderBy"`
	Filter   string `json:"filter"`
}

// type User struct {
// 	ID       string `db:"id" json:"id"`
// 	UserName string `db:"username" json:"username"`
// 	Email    string `db:"email" json:"email"`
// 	IsAdmin  *bool  `db:"is_admin" json:"isAdmin"`
// 	Rate     int    `db:"rate" json:"rate"`
// 	Deleted  bool   `db:"deleted" json:"-"`
// }

type User struct {
	ID          int     `db:"id" json:"id"`
	UserName    string  `db:"user_name" json:"userName"`
	FirstName   string  `db:"first_name" json:"firstName"`
	LastName    string  `db:"last_name" json:"lastName"`
	Email       string  `db:"email" json:"email"`
	CacUid      *string `db:"edipi" json:"cacUid"`
	RoleID      int     `db:"role_id" json:"roleId"`
	OfficeID    int     `db:"office_id" json:"officeId"`
	Role        string  `db:"description" json:"role"`
	OfficeCode  string  `db:"code" json:"officeCode"`
	ProjectCode string  `db:"project_code" json:"projectCode"`
}

type UserRoleOffice struct {
	ID          int    `db:"id" json:"id"`
	UserID      int    `db:"user_id" json:"userId"`
	RoleID      int    `db:"role_id" json:"roleId"`
	OfficeID    int    `db:"office_id" json:"officeId"`
	Role        string `db:"description" json:"role"`
	OfficeCode  string `db:"code" json:"officeCode"`
	ProjectCode string `db:"project_code" json:"projectCode"`
}

type KeyCloakResponse struct {
	AccessToken string `json:"access_token"`
	Expires     int    `json:"expires_in"`
	ClientId    string `json:"id"`
}

type KeyCloakUser struct {
	Username  string       `json:"username"`
	FirstName string       `json:"firstName"`
	LastName  string       `json:"lastName"`
	UserID    string       `json:"id"`
	Role      KeyCloakRole `json:"role"`
}

type KeyCloakRole struct {
	RoleName string `json:"roleName"`
	RoleId   string `json:"id"`
}

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}
