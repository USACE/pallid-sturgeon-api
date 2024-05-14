package stores

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/stdlib"

	_ "github.com/godror/godror"

	"github.com/USACE/pallid_sturgeon_api/server/config"
)

func InitStores(appConfig *config.AppConfig) (*PallidSturgeonStore, error) {
	connectString := fmt.Sprintf("%s:%s/%s", appConfig.Dbhost, appConfig.Dbport, appConfig.Dbname)
	db, err := sqlx.Connect(
		"godror",
		"user="+appConfig.Dbuser+" password="+appConfig.Dbpass+" connectString="+connectString+" poolMaxSessions=100 poolSessionMaxLifetime=2m0s",
	)
	
	if err != nil {
		log.Printf("[InitStores] m=GetDb,msg=connection has failed: %s", err)
		return nil, err
	} else {
		db.SetMaxIdleConns(0)	
	}

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
