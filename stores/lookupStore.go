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

func (s *LookupStore) GetBendSelections() ([]models.BendSelection, error) {
	query := `
        SELECT bs_id, bend_selection_code, bend_selection_description
		FROM bend_selection_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.BendSelection{}

	for rows.Next() {
		var i models.BendSelection
		if err := rows.Scan(&i.BsId, &i.BendSelectionCode, &i.BendSelectionDesc); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetGearCodes() ([]models.GearCode, error) {
	query := `
        SELECT gear_id, gear, gear_code, gear_type, gear_description
		FROM gear_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.GearCode{}

	for rows.Next() {
		var i models.GearCode
		if err := rows.Scan(&i.GearId, &i.Gear, &i.GearCode, &i.GearType, &i.GearDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetFilteredGearCodes() ([]models.FilteredGearCode, error) {
	query := `
        SELECT sgo_id, field_office_code, season_code, gear_code, gear, project_code
		FROM season_gear_office_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.FilteredGearCode{}

	for rows.Next() {
		var i models.FilteredGearCode
		if err := rows.Scan(&i.SgoId, &i.FieldOfficeCode, &i.SeasonCode, &i.GearCode, &i.Gear, &i.ProjectCode); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetGearTypes() ([]models.GearType, error) {
	query := `
        SELECT gt_id, gear_type_code, gear_type_description
		FROM gear_type_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.GearType{}

	for rows.Next() {
		var i models.GearType
		if err := rows.Scan(&i.GtId, &i.GearTypeCode, &i.GearTypeDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetMacros() ([]models.Macro, error) {
	query := `
        SELECT mh_id, habitat_code, habitat_description
		FROM macrohabitat_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.Macro{}

	for rows.Next() {
		var i models.Macro
		if err := rows.Scan(&i.MhId, &i.HabitatCode, &i.HabitatDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetMesos() ([]models.MesoLk, error) {
	query := `
        SELECT mh_id, mesohabitat_code, mesohabitat_description
		FROM mesohabitat_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.MesoLk{}

	for rows.Next() {
		var i models.MesoLk
		if err := rows.Scan(&i.MhId, &i.MesoHabitatCode, &i.MesoHabitatDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetMacroMesos() ([]models.MacroMeso, error) {
	query := `
        SELECT mm_id, macrohabitat_code, mesohabitat_code
		FROM macro_meso_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.MacroMeso{}

	for rows.Next() {
		var i models.MacroMeso
		if err := rows.Scan(&i.MmId, &i.MacroHabitatCode, &i.MesoHabitatCode); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetMicroHabitats() ([]models.MicroHabitat, error) {
	query := `
        SELECT mh_id, micro_structure, micro_structure_code, structure_flow, structure_flow_code, structure_mod, structure_mod_code
		FROM micro_habitat_desc_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.MicroHabitat{}

	for rows.Next() {
		var i models.MicroHabitat
		if err := rows.Scan(&i.MhId, &i.MicroStructure, &i.MicroStructureCode, &i.StructureFlow, &i.StructureFlowCode, &i.StructureMod, &i.StructureModCode); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetU7() ([]models.U7, error) {
	query := `
        SELECT code, description
		FROM useven_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.U7{}

	for rows.Next() {
		var i models.U7
		if err := rows.Scan(&i.U7Code, &i.U7Description); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}
