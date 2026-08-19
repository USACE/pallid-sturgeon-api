package stores

import (
	"fmt"
	"log"

	"github.com/USACE/pallid_sturgeon_api/server/config"
	"github.com/USACE/pallid_sturgeon_api/server/models"
	"github.com/jmoiron/sqlx"
)

type LookupStore struct {
	db *sqlx.DB
}

func NewLookupStore(db *sqlx.DB) *LookupStore {
	return &LookupStore{db}
}

func InitLookupStore(appConfig *config.AppConfig) (*LookupStore, error) {
	connectString := fmt.Sprintf("%s:%s/%s", appConfig.Dbhost, appConfig.Dbport, appConfig.Dbname)
	db, err := sqlx.Connect(
		"godror",
		"user="+appConfig.Dbuser+" password="+appConfig.Dbpass+" connectString="+connectString+" poolMaxSessions=100 poolSessionMaxLifetime=2m0s",
	)
	db.SetMaxIdleConns(0)
	if err != nil {
		log.Printf("[InitAuthStore] m=GetDb,msg=connection has failed: %s", err)
		return nil, err
	}

	ss := LookupStore{
		db: db,
	}

	return &ss, nil
}

// GLOBAL LOOK UP QUERIES

func (s *LookupStore) GetFrequencyIds() ([]models.FrequencyId, error) {
	query := `
        SELECT frequency_id_code, frequency_id_description FROM frequency_id_lk
		WHERE active_flag_tf = 'T' ORDER BY sort_order asc
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.FrequencyId{}

	for rows.Next() {
		var i models.FrequencyId
		if err := rows.Scan(&i.FrequencyIdCode, &i.FrequencyIdDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetScuteLocations() ([]models.ScuteLocation, error) {
	query := `
        SELECT scute_location_code, scute_location_description
		FROM scute_location_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.ScuteLocation{}

	for rows.Next() {
		var i models.ScuteLocation
		if err := rows.Scan(&i.ScuteLocationCode, &i.ScuteLocationDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}