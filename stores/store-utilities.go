package stores

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/stdlib"

	_ "github.com/godror/godror"

	"github.com/USACE/pallid_sturgeon_api/server/config"
)

func InitStores(appConfig *config.AppConfig) (*PallidSturgeonStore, error) {
	dburl := fmt.Sprintf("user=\"%s\" password=\"%s\" connectString=\"%s:%s/%s\" poolMaxSessions=50 poolSessionTimeout=42s",
		appConfig.Dbuser, appConfig.Dbpass, appConfig.Dbhost, appConfig.Dbport, appConfig.Dbname)
	db, err := sql.Open("godror", dburl)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//defer db.Close()

	ss := PallidSturgeonStore{
		db:     db,
		config: appConfig,
	}

	return &ss, nil
}

type TransactionFunction func(*sqlx.Tx)

/*
Transaction Wrapper.
DB Calls within the transaction should panic on fail.  i.e. use MustExec vs Exec.
*/

// func transaction(db *sqlx.DB, fn TransactionFunction) error {
// 	var err error
// 	tx, err := db.Beginx()
// 	if err != nil {
// 		log.Printf("Unable to start transaction: %s\n", err)
// 		return err
// 	}
// 	defer func() {
// 		if r := recover(); r != nil {
// 			log.Print(r)
// 			err = tx.Rollback()
// 			if err != nil {
// 				log.Printf("Unable to rollback from transaction: %s", err)
// 			}
// 		} else {
// 			err = tx.Commit()
// 			if err != nil {
// 				log.Printf("Unable to commit transaction: %s", err)
// 			}
// 		}
// 	}()
// 	fn(tx)
// 	return err
// }
