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

func (s *LookupStore) GetYears() ([]models.YearLK, error) {
	query := `
        SELECT year_id, year
		FROM year_lk ORDER BY year desc
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.YearLK{}

	for rows.Next() {
		var i models.YearLK
		if err := rows.Scan(&i.YearId, &i.Year); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetFieldOffices() ([]models.FieldOfficeLK, error) {
	query := `
        SELECT fo_id, field_office_code, field_office_description, state
		FROM field_office_lk ORDER BY fo_id desc
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.FieldOfficeLK{}

	for rows.Next() {
		var i models.FieldOfficeLK
		if err := rows.Scan(&i.FoId, &i.FoCode, &i.FoDescription, &i.State); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetFieldOfficeSegment() ([]models.FieldOfficeSegment, error) {
	query := `
        SELECT fos_id, field_office_code, segment_code, project_code
		FROM field_office_segment_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.FieldOfficeSegment{}

	for rows.Next() {
		var i models.FieldOfficeSegment
		if err := rows.Scan(&i.FosId, &i.FieldOfficeCode, &i.SegmentCode, &i.ProjectCode); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetProjects() ([]models.ProjectLK, error) {
	query := `
        SELECT project_code, project_description
		FROM project_lk WHERE active_flag_tf = 'T'
		ORDER BY project_code asc
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.ProjectLK{}

	for rows.Next() {
		var i models.ProjectLK
		if err := rows.Scan(&i.ProjectCode, &i.ProjectDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetSegments() ([]models.SegmentLK, error) {
	query := `
        SELECT s_id, segment_code, segment_description
		FROM segment_lk ORDER BY segment_code asc
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.SegmentLK{}

	for rows.Next() {
		var i models.SegmentLK
		if err := rows.Scan(&i.SegmentId, &i.SegmentCode, &i.SegmentDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetSeasons() ([]models.SeasonLK, error) {
	query := `
        SELECT s_id, season_code, season_description
		FROM season_lk WHERE active_flag_tf = 'T'
		ORDER BY s_id asc
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.SeasonLK{}

	for rows.Next() {
		var i models.SeasonLK
		if err := rows.Scan(&i.SeasonId, &i.SeasonCode, &i.SeasonDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetSampleUnitTypes() ([]models.SampleUnitTypeLK, error) {
	query := `
        SELECT sample_unit_type_code, sample_unit_type_description
		FROM sample_unit_type_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.SampleUnitTypeLK{}

	for rows.Next() {
		var i models.SampleUnitTypeLK
		if err := rows.Scan(&i.SutCode, &i.SutDescription ); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
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

func (s *LookupStore) GetBendRiverMile() ([]models.BendRiverMile, error) {
	query := `
        SELECT brm_id, b_segment, bend_num, b_desc, state, upper_river_mile, lower_river_mile
		FROM bend_river_mile_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.BendRiverMile{}

	for rows.Next() {
		var i models.BendRiverMile
		if err := rows.Scan(&i.BrmId, &i.Segment, &i.Bend, &i.BendDescription, &i.State, &i.UpperRiverMile, &i.LowerRiverMile); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetChutes() ([]models.Chute, error) {
	query := `
        SELECT chute_id, segment_id, chute_code, chute_desc, upper_river_mile
		FROM chute_lk ORDER BY chute_desc asc
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.Chute{}

	for rows.Next() {
		var i models.Chute
		if err := rows.Scan(&i.ChuteId, &i.Segment, &i.ChuteCode, &i.ChuteDescription, &i.UpperRiverMile); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetReach() ([]models.Reach, error) {
	query := `
        SELECT reach_id, segment_id, reach_code, reach_desc, upper_river_mile
		FROM reach_lk ORDER BY reach_desc asc
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.Reach{}

	for rows.Next() {
		var i models.Reach
		if err := rows.Scan(&i.ReachId, &i.Segment, &i.ReachCode, &i.ReachDescription, &i.UpperRiverMile); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetRecaptureData() ([]models.Recapture, error) {
	query := `
        SELECT species, pit_tag
		FROM recapture_data
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.Recapture{}

	for rows.Next() {
		var i models.Recapture
		if err := rows.Scan(&i.Species, &i.PitTag); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}


// MISSOURI RIVER LOOK UP QUERIES

func (s *LookupStore) GetGearCodes() ([]models.GearCode, error) {
	query := `
        SELECT gear_id, gear, gear_code, gear_type, gear_description, deploymenttype
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
		if err := rows.Scan(&i.GearId, &i.Gear, &i.GearCode, &i.GearType, &i.GearDescription, &i.DeploymentType); err != nil {
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

func (s *LookupStore) GetMicroStructures() ([]models.MicroStructure, error) {
	query := `
        SELECT ms_id, code, description
		FROM micro_structure_lk WHERE active_flag_tf = 'T'
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.MicroStructure{}

	for rows.Next() {
		var i models.MicroStructure
		if err := rows.Scan(&i.MsId, &i.MicroStructureCode, &i.MicroStructureDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetStructureFlow() ([]models.StructureFlowLK, error) {
	query := `
        SELECT sf_id, code, description
		FROM structure_flow_lk WHERE active_flag_tf = 'T'
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.StructureFlowLK{}

	for rows.Next() {
		var i models.StructureFlowLK
		if err := rows.Scan(&i.SfId, &i.StructureFlowCode, &i.StructureFlowDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetStructureMod() ([]models.StructureModLK, error) {
	query := `
        SELECT sm_id, code, description
		FROM structure_mod_lk WHERE active_flag_tf = 'T'
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.StructureModLK{}

	for rows.Next() {
		var i models.StructureModLK
		if err := rows.Scan(&i.SmId, &i.StructureModCode, &i.StructureModDescription); err != nil {
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

func (s *LookupStore) GetU6() ([]models.U6, error) {
	query := `
        SELECT u6_id, code, description
		FROM u6_lk WHERE active_flag_tf = 'T'
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.U6{}

	for rows.Next() {
		var i models.U6
		if err := rows.Scan(&i.U6Id, &i.U6Code, &i.U6Description); err != nil {
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

func (s *LookupStore) GetEstimations() ([]models.Estimation, error) {
	query := `
        SELECT estimation_code, estimation_description
		FROM cobble_organic_est_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.Estimation{}

	for rows.Next() {
		var i models.Estimation
		if err := rows.Scan(&i.EstCode, &i.EstDesc); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetMicroSetSite() ([]models.MicroSetSite, error) {
	query := `
        SELECT ms_id, structure_code, micro_structure, set_site_1_code, set_site_1, set_site_two_code, set_site_two
		FROM micro_set_site_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.MicroSetSite{}

	for rows.Next() {
		var i models.MicroSetSite
		if err := rows.Scan(&i.MsId, &i.MicroStructureCode, &i.MicroStructureDesc, &i.Ss1Code, &i.Ss1Description, &i.Ss2Code, &i.Ss2Description); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetSetSite1() ([]models.SetSite1LK, error) {
	query := `
        SELECT ss1_id, code, description
		FROM set_site_1_lk WHERE active_flag_tf = 'T'
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.SetSite1LK{}

	for rows.Next() {
		var i models.SetSite1LK
		if err := rows.Scan(&i.Ss1Id, &i.Ss1Code, &i.Ss1Description); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetSetSite2() ([]models.SetSite2LK, error) {
	query := `
        SELECT ss2_id, code, description
		FROM set_site_2_lk WHERE active_flag_tf = 'T'
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.SetSite2LK{}

	for rows.Next() {
		var i models.SetSite2LK
		if err := rows.Scan(&i.Ss2Id, &i.Ss2Code, &i.Ss2Description); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetSetSite3() ([]models.SetSite3, error) {
	query := `
        SELECT set_site_3_code, set_site_3
		FROM micro_set_site3_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.SetSite3{}

	for rows.Next() {
		var i models.SetSite3
		if err := rows.Scan(&i.SsCode, &i.SsDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetSubsampleTypes() ([]models.SubsampleType, error) {
	query := `
        SELECT st_id, code, description
		FROM subsample_type_lk WHERE active_flag_tf = 'T'
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.SubsampleType{}

	for rows.Next() {
		var i models.SubsampleType
		if err := rows.Scan(&i.StId, &i.StCode, &i.StDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

// FISH LOOK UP QUERIES

func (s *LookupStore) GetFishCodes() ([]models.FishCode, error) {
	query := `
        SELECT fish_id, common_name, scientific_name, alpha_code, numeric_codes, numeric_codes_txt
		FROM fish_code_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.FishCode{}

	for rows.Next() {
		var i models.FishCode
		if err := rows.Scan(&i.FishId, &i.CommonName, &i.ScientificName, &i.AlphaCode, &i.NumericCodes, &i.NumericCodesText); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetFishStructures() ([]models.FishStructure, error) {
	query := `
        SELECT code, description
		FROM fish_structure_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.FishStructure{}

	for rows.Next() {
		var i models.FishStructure
		if err := rows.Scan(&i.FsCode, &i.FsDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetFloyTagPrefixes() ([]models.FloyTagPrefix, error) {
	query := `
        SELECT floy_id, tag_prefix_code, tag_prefix_description
		FROM floy_tag_prefix_code_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.FloyTagPrefix{}

	for rows.Next() {
		var i models.FloyTagPrefix
		if err := rows.Scan(&i.FtpId, &i.FtpCode, &i.FtpDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetLengthTypes() ([]models.LengthType, error) {
	query := `
        SELECT lt_id, code, description
		FROM length_type_lk WHERE active_flag_tf = 'T'
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.LengthType{}

	for rows.Next() {
		var i models.LengthType
		if err := rows.Scan(&i.LtId, &i.LtCode, &i.LtDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetMarkRecapture() ([]models.MarkRecapture, error) {
	query := `
        SELECT mr_id, mark_recapture_code, mark_recapture_description
		FROM mark_recapture_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.MarkRecapture{}

	for rows.Next() {
		var i models.MarkRecapture
		if err := rows.Scan(&i.MrId, &i.MrCode, &i.MrDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetFrequencyId() ([]models.FrequencyId, error) {
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

func (s *LookupStore) GetSpawnBehavior() ([]models.SpawnBehavior, error) {
	query := `
        SELECT spawn_code, spawn_description
		FROM spawn_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.SpawnBehavior{}

	for rows.Next() {
		var i models.SpawnBehavior
		if err := rows.Scan(&i.SpawnCode, &i.SpawnDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetPositionConfidence() ([]models.PositionConfidence, error) {
	query := `
        SELECT position_confidence_code, position_confidence_description
		FROM position_confidence_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.PositionConfidence{}

	for rows.Next() {
		var i models.PositionConfidence
		if err := rows.Scan(&i.PositionCode, &i.PositionDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (s *LookupStore) GetSearchType() ([]models.SearchType, error) {
	query := `
        SELECT search_type_code, search_type_description
		FROM search_type_lk
    `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []models.SearchType{}

	for rows.Next() {
		var i models.SearchType
		if err := rows.Scan(&i.SearchTypeCode, &i.SearchTypeDescription); err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}
