package stores

import (
	"database/sql"
	"fmt"
	"log"

	"di2e.net/cwbi/pallid_sturgeon_api/server/config"
	"di2e.net/cwbi/pallid_sturgeon_api/server/models"
	"github.com/jmoiron/sqlx"
)

type AuthStore struct {
	db *sqlx.DB
}

var userSql = "select * from users where user_id=$1"

// var userSql = `select id,username,email,rate,
// 				(select bool_or(is_admin)
// 					from org_members
// 					where user_id=$1) as is_admin
// 				from wpt_user
// 				where id=$1 and deleted=false`

var insertUserSql = "insert into users (user_id,user_name,email) values ($1,$2,$3)"

func InitAuthStore(appConfig *config.AppConfig) (*AuthStore, error) {
	dburl := fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=disable",
		appConfig.Dbuser, appConfig.Dbpass, appConfig.Dbhost, appConfig.Dbport, appConfig.Dbname)
	con, err := sqlx.Connect("pgx", dburl)
	if err != nil {
		log.Printf("Unable to connect to the authentication store: %s", err)
		return nil, err
	}
	con.SetMaxOpenConns(10)
	adb := AuthStore{
		db: con,
	}
	return &adb, nil
}

func (auth *AuthStore) GetUserFromJwt(jwtClaims models.JwtClaim) (models.User, error) {
	user := models.User{}
	err := auth.db.Get(&user, userSql, jwtClaims.Sub)
	if err != nil {
		if err == sql.ErrNoRows {
			user = models.User{
				UserID:   jwtClaims.Sub,
				UserName: jwtClaims.Name,
				Email:    jwtClaims.Email,
				Deleted:  false,
			}
			err := auth.AddUser(user)
			return user, err

		} else {
			return user, err
		}
	}
	return user, err
}

func (auth *AuthStore) AddUser(user models.User) error {
	_, err := auth.db.Exec(insertUserSql, user.UserID, user.UserName, user.Email)
	return err
}
