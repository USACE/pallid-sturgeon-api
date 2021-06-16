package models

type JwtClaim struct {
	Sub   string
	Name  string
	Email string
	Roles []interface{}
}

type SearchParams struct {
	Page        int    `json:"page"`
	PageSize    int    `json:"size"`
	OrderBy     string `json:"orderBy"`
	Filter      string `json:"filter"`
	PhaseType   string `json:"phaseType"`
	PhaseStatus string `json:"phaseStatus"`
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
	UserID   string `db:"user_id" json:"userId"`
	UserName string `db:"user_name" json:"userName"`
	Email    string `db:"email" json:"email"`
	Deleted  bool   `db:"deleted" json:"-"`
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
