package stores

import (
	"fmt"
	"log"

	//"net/smtp"

	"github.com/USACE/pallid_sturgeon_api/server/config"
	"github.com/USACE/pallid_sturgeon_api/server/models"
	"github.com/jmoiron/sqlx"
)

type AuthStore struct {
	db *sqlx.DB
	//config *config.AppConfig
}

var userSql = "select id, edipi, username, email, first_name,last_name from users_t where email=:1"

var userByIdSql = "select id, edipi, username, email, first_name,last_name from users_t where id=:1"

// var userSql = `select id,username,email,rate,
// 				(select bool_or(is_admin)
// 					from org_members
// 					where user_id=$1) as is_admin
// 				from wpt_user
// 				where id=$1 and deleted=false`

var insertUserSql = "insert into users_t (username,email,first_name,last_name,edipi) values (:1,:2,:3,:4,:5)"

var getUsersSql = `select u.id, u.username, u.first_name, u.last_name, u.email, uro.role_id, r.description, uro.office_id, f.field_office_code, uro.project_code from users_t u 
	inner join user_role_office_lk uro on uro.user_id = u.id 
	inner join role_lk r on r.id = uro.role_id 
	inner join field_office_lk f on f.fo_id = uro.office_id order by u.last_name`

var getUsersByRoleTypeSql = `select u.id, u.username, u.first_name, u.last_name, u.email, r.description from users_t u
						inner join user_role_office_lk uro on uro.user_id = u.id
						inner join role_lk r on r.id = uro.role_id
						where r.description = :1`

var getUserAccessRequestSql = "select id, username, first_name, last_name, email from users_t where id not in (select user_id from user_role_office_lk) order by last_name"

var insertUserRoleOfficeSql = "insert into user_role_office_lk (id,user_id,role_id,office_id,project_code) values (user_role_office_seq.nextval,:1,:2,:3,:4)"

var updateUserRoleOfficesSql = `update user_role_office_lk set role_id = :2, office_id = :3, project_code = :4 where user_id = :1`

var getUserRoleOfficeSql = `select uro.id, uro.user_id, uro.role_id, uro.office_id, r.description, f.FIELD_OFFICE_CODE, project_code from user_role_office_lk uro 
							inner join users_t u on u.id = uro.user_id
							inner join role_lk r on r.id = uro.role_id
							inner join field_office_lk f on f.fo_id = uro.office_id
							where u.email = :1`

func InitAuthStore(appConfig *config.AppConfig) (*AuthStore, error) {
	connectString := fmt.Sprintf("%s:%s/%s", appConfig.Dbhost, appConfig.Dbport, appConfig.Dbname)
	db, err := sqlx.Connect(
		"godror",
		"user="+appConfig.Dbuser+" password="+appConfig.Dbpass+" connectString="+connectString+" poolMaxSessions=100 poolSessionMaxLifetime=2m0s",
	)
	// db.SetMaxIdleConns(2)
	// db.SetConnMaxLifetime(2 * time.Minute)
	// db.SetMaxOpenConns(100)
	if err != nil {
		log.Printf("[InitAuthStore] m=GetDb,msg=connection has failed: %s", err)
		return nil, err
	}

	ss := AuthStore{
		db: db,
		//config: appConfig,
	}

	return &ss, nil
}

func (auth *AuthStore) GetUserFromJwt(jwtClaims models.JwtClaim) (models.User, error) {
	user := models.User{}

	countQuery, err := auth.db.Prepare(userSql)
	if err != nil {
		return user, err
	}

	countrows, err := countQuery.Query(jwtClaims.Email)
	if err != nil {
		return user, err
	}
	count := 0
	for countrows.Next() {
		err = countrows.Scan(&user.ID, &user.CacUid, &user.UserName, &user.Email, &user.FirstName, &user.LastName)
		if err != nil {
			return user, err
		}
		count += 1
	}

	if count == 0 {
		user = models.User{
			UserName:  jwtClaims.Email,
			Email:     jwtClaims.Email,
			FirstName: jwtClaims.FirstName,
			LastName:  jwtClaims.LastName,
			CacUid:    jwtClaims.CacUid,
		}
		err := auth.AddUser(user)
		return user, err
	}
	defer countrows.Close()
	return user, err
}

func (auth *AuthStore) AddUser(user models.User) error {
	_, err := auth.db.Exec(insertUserSql, user.UserName, user.Email, user.FirstName, user.LastName, user.CacUid)
	return err
}

func (auth *AuthStore) GetUserAccessRequests() ([]models.User, error) {
	users := []models.User{}

	rows, err := auth.db.Query(getUserAccessRequestSql)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName, &user.Email)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	defer rows.Close()

	return users, err
}

func (auth *AuthStore) GetUsers() ([]models.User, error) {
	users := []models.User{}

	rows, err := auth.db.Query(getUsersSql)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.RoleID, &user.Role, &user.OfficeID, &user.OfficeCode, &user.ProjectCode)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	defer rows.Close()

	return users, err
}

func (auth *AuthStore) UpdateUserRoleOffice(userRoleOffice models.UserRoleOffice) error {
	_, err := auth.db.Exec(updateUserRoleOfficesSql, userRoleOffice.RoleID, userRoleOffice.OfficeID, userRoleOffice.ProjectCode, userRoleOffice.UserID)

	return err
}

func (auth *AuthStore) GetUserById(id int) (models.User, error) {
	user := models.User{}

	countQuery, err := auth.db.Prepare(userByIdSql)
	if err != nil {
		return user, err
	}

	countrows, err := countQuery.Query(id)
	if err != nil {
		return user, err
	}
	count := 0
	for countrows.Next() {
		err = countrows.Scan(&user.ID, &user.CacUid, &user.UserName, &user.Email, &user.FirstName, &user.LastName)
		if err != nil {
			return user, err
		}
		count += 1
	}

	defer countrows.Close()
	return user, err
}

func (auth *AuthStore) GetUsersByRoleType(roleType string) ([]models.User, error) {
	users := []models.User{}

	selectQuery, err := auth.db.Prepare(getUsersByRoleTypeSql)
	if err != nil {
		return users, err
	}

	rows, err := selectQuery.Query(roleType)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.Role)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	defer rows.Close()

	return users, err
}

func (auth *AuthStore) AddUserRoleOffice(userRoleOffice models.UserRoleOffice) error {
	_, err := auth.db.Exec(insertUserRoleOfficeSql, userRoleOffice.UserID, userRoleOffice.RoleID, userRoleOffice.OfficeID, userRoleOffice.ProjectCode)
	// if err == nil {
	// 	message := []byte("Your role request has been approved.")

	// 	user, userErr := auth.GetUserById(userRoleOffice.UserID)
	// 	if userErr != nil {
	// 		log.Print("Unable to send email.", userErr)
	// 	}
	// 	to := []string{
	// 		user.Email,
	// 	}

	// 	from := auth.config.EmailFrom
	// 	emailErr := auth.SendEmail(message, to, from)
	// 	if emailErr != nil {
	// 		log.Print("Unable to send email.", emailErr)
	// 	}
	// }

	return err
}

func (auth *AuthStore) GetUserRoleOffice(email string) (models.UserRoleOffice, error) {
	userRoleOffice := models.UserRoleOffice{}
	selectQuery, err := auth.db.Prepare(getUserRoleOfficeSql)
	if err != nil {
		return userRoleOffice, err
	}

	rows, err := selectQuery.Query(email)
	if err != nil {
		return userRoleOffice, err
	}

	for rows.Next() {
		err = rows.Scan(&userRoleOffice.ID, &userRoleOffice.UserID, &userRoleOffice.RoleID, &userRoleOffice.OfficeID, &userRoleOffice.Role, &userRoleOffice.OfficeCode, &userRoleOffice.ProjectCode)
		if err != nil {
			return userRoleOffice, err
		}
	}
	defer rows.Close()

	// if userRoleOffice.OfficeCode == "" {
	// 	message := []byte("There is a new user role request. Please login to appove or deny the request.")
	// 	users, adminUserLoadErr := auth.GetUsersByRoleType("ADMINISTRATOR")
	// 	if adminUserLoadErr != nil {
	// 		log.Print("Unable to send email.", adminUserLoadErr)
	// 	}

	// 	to := make([]string, 0)
	// 	for _, user := range users {
	// 		to = append(to, user.Email)
	// 	}

	// 	from := auth.config.EmailFrom
	// 	emailErr := auth.SendEmail(message, to, from)
	// 	if emailErr != nil {
	// 		log.Print("Unable to send email.", emailErr)
	// 	}
	// }

	return userRoleOffice, err
}

// func (s *AuthStore) SendEmail(message []byte, to []string, from string) error {
// 	// Authentication.
// 	auth := smtp.PlainAuth("", from, s.config.EmailPassword, s.config.SmtpHost)

// 	// Sending email.
// 	err := smtp.SendMail(s.config.SmtpHost+":"+s.config.SmtpPort, auth, from, to, message)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	fmt.Println("Email Sent Successfully!")
// 	return nil
// }
