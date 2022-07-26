package stores

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/USACE/pallid_sturgeon_api/server/config"
	"github.com/USACE/pallid_sturgeon_api/server/models"
	"github.com/godror/godror"
	"github.com/jmoiron/sqlx"
)

type PallidSturgeonStore struct {
	db     *sqlx.DB
	config *config.AppConfig
}

var getUserSql = `select u.id, u.username, u.first_name, u.last_name, u.email, r.description, f.FIELD_OFFICE_CODE from users_t u
							inner join user_role_office_lk uro on uro.user_id = u.id
							inner join role_lk r on r.id = uro.role_id
							inner join field_office_lk f on f.fo_id = uro.office_id
				    where email = :1`

func (s *PallidSturgeonStore) GetUser(email string) (models.User, error) {
	user := models.User{}
	selectQuery, err := s.db.Prepare(getUserSql)
	if err != nil {
		return user, err
	}

	rows, err := selectQuery.Query(email)
	if err != nil {
		return user, err
	}

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.OfficeCode)
		if err != nil {
			return user, err
		}
	}
	defer rows.Close()

	return user, err
}

var getProjectsSql = `select distinct p.* from project_lk p
						join fieldoffice_segment_v v
						on v.PROJECT_CODE = p.project_code
						and v.FIELD_OFFICE_CODE = :1
						order by p.project_code`

func (s *PallidSturgeonStore) GetProjects(fieldOfficeCode string) ([]models.Project, error) {
	projects := []models.Project{}

	selectQuery, err := s.db.Prepare(getProjectsSql)
	if err != nil {
		return projects, err
	}

	rows, err := selectQuery.Query(fieldOfficeCode)
	if err != nil {
		return projects, err
	}
	defer rows.Close()

	for rows.Next() {
		project := models.Project{}
		err = rows.Scan(&project.Code, &project.Description)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, err
}

func (s *PallidSturgeonStore) GetRoles() ([]models.Role, error) {
	rows, err := s.db.Query("select * from role_lk order by id")

	roles := []models.Role{}
	if err != nil {
		return roles, err
	}

	for rows.Next() {
		role := models.Role{}
		err = rows.Scan(&role.ID, &role.Description)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	defer rows.Close()

	return roles, err
}

func (s *PallidSturgeonStore) GetFieldOffices() ([]models.FieldOffice, error) {
	rows, err := s.db.Query("select FO_ID, FIELD_OFFICE_CODE, FIELD_OFFICE_DESCRIPTION, STATE from field_office_lk order by FO_ID")

	fieldOffices := []models.FieldOffice{}
	if err != nil {
		return fieldOffices, err
	}
	defer rows.Close()

	for rows.Next() {
		fieldOffice := models.FieldOffice{}
		err = rows.Scan(&fieldOffice.ID, &fieldOffice.Code, &fieldOffice.Description, &fieldOffice.State)
		if err != nil {
			return nil, err
		}
		fieldOffices = append(fieldOffices, fieldOffice)
	}

	return fieldOffices, err
}

func (s *PallidSturgeonStore) GetSampleMethods() ([]models.SampleMethod, error) {
	rows, err := s.db.Query("select * from sample_type_lk order by SAMPLE_TYPE_CODE")

	sampleMethods := []models.SampleMethod{}
	if err != nil {
		return sampleMethods, err
	}
	defer rows.Close()

	for rows.Next() {
		sampleMethod := models.SampleMethod{}
		err = rows.Scan(&sampleMethod.Code, &sampleMethod.Description)
		if err != nil {
			return nil, err
		}
		sampleMethods = append(sampleMethods, sampleMethod)
	}

	return sampleMethods, err
}

func (s *PallidSturgeonStore) GetSampleUnitTypes() ([]models.SampleUnitType, error) {
	rows, err := s.db.Query("select * from sample_unit_type_lk order by SAMPLE_UNIT_TYPE_CODE")

	sampleUnitTypes := []models.SampleUnitType{}
	if err != nil {
		return sampleUnitTypes, err
	}
	defer rows.Close()

	for rows.Next() {
		sampleUnitType := models.SampleUnitType{}
		err = rows.Scan(&sampleUnitType.Code, &sampleUnitType.Description)
		if err != nil {
			return nil, err
		}
		sampleUnitTypes = append(sampleUnitTypes, sampleUnitType)
	}

	return sampleUnitTypes, err
}

func (s *PallidSturgeonStore) GetSeasons() ([]models.Season, error) {
	rows, err := s.db.Query("select s_id, season_code, season_description, field_app, PROJECT_CODE from season_lk order by s_id")

	seasons := []models.Season{}
	if err != nil {
		return seasons, err
	}
	defer rows.Close()

	for rows.Next() {
		season := models.Season{}
		err = rows.Scan(&season.ID, &season.Code, &season.Description, &season.FieldAppFlag, &season.ProjectCode)
		if err != nil {
			return nil, err
		}
		seasons = append(seasons, season)
	}

	return seasons, err
}

var getSegmentsSql = `select distinct s.s_id, s.segment_code, s.segment_description, s.segment_type, s.river, s.upper_river_mile, s.lower_river_mile, s.rpma from segment_lk s
						join fieldoffice_segment_v v
						on v.SEGMENT_CODE = s.segment_code
						and v.FIELD_OFFICE_CODE = :1
						order by s.s_id`

func (s *PallidSturgeonStore) GetSegments(fieldOfficeCode string) ([]models.Segment, error) {
	segments := []models.Segment{}

	selectQuery, err := s.db.Prepare(getSegmentsSql)
	if err != nil {
		return segments, err
	}

	rows, err := selectQuery.Query(fieldOfficeCode)
	if err != nil {
		return segments, err
	}
	defer rows.Close()

	for rows.Next() {
		segment := models.Segment{}
		err = rows.Scan(&segment.ID, &segment.Code, &segment.Description, &segment.Type, &segment.RiverCode, &segment.UpperRiverMile, &segment.LowerRiverMile, &segment.Rpma)
		if err != nil {
			return nil, err
		}
		segments = append(segments, segment)
	}

	return segments, err
}

func (s *PallidSturgeonStore) GetBends() ([]models.Bend, error) {
	rows, err := s.db.Query("select BRM_ID, BEND_NUM, B_DESC, B_SEGMENT, upper_river_mile, lower_river_mile, state from bend_river_mile_lk order by BRM_ID")

	bends := []models.Bend{}
	if err != nil {
		return bends, err
	}
	defer rows.Close()

	for rows.Next() {
		bend := models.Bend{}
		err = rows.Scan(&bend.ID, &bend.BendNumber, &bend.Description, &bend.SegmentCode, &bend.UpperRiverMile, &bend.LowerRiverMile, &bend.State)
		if err != nil {
			return nil, err
		}
		bends = append(bends, bend)
	}

	return bends, err
}

var siteDataEntriesSql = `SELECT brm_id, site_id, site_fid, year, FIELDOFFICE, PROJECT_ID,
SEGMENT_ID, SEASON, SAMPLE_UNIT_TYPE, BENDRN, edit_initials, uploaded_by fROM ds_sites where year = :1 and FIELDOFFICE = :2`

var siteDataEntriesCountSql = `SELECT count(*) FROM ds_sites where year = :1 and FIELDOFFICE = :2`

func (s *PallidSturgeonStore) GetSiteDataEntries(year string, fieldOfficeCode string, projectCode string, segmentCode string, seasonCode string, bendrn string, queryParams models.SearchParams) (models.SitesWithCount, error) {
	siteDataEntryWithCount := models.SitesWithCount{}
	query := siteDataEntriesSql
	queryWithCount := siteDataEntriesCountSql

	if projectCode != "" {
		query = query + fmt.Sprintf(" and PROJECT_ID = '%s' ", projectCode)
		queryWithCount = queryWithCount + fmt.Sprintf(" and PROJECT_ID = '%s' ", projectCode)
	}

	if segmentCode != "" {
		query = query + fmt.Sprintf(" and SEGMENT_ID = '%s' ", segmentCode)
		queryWithCount = queryWithCount + fmt.Sprintf(" and SEGMENT_ID = '%s' ", segmentCode)
	}

	if seasonCode != "" {
		query = query + fmt.Sprintf(" and SEASON = '%s' ", seasonCode)
		queryWithCount = queryWithCount + fmt.Sprintf(" and SEASON = '%s' ", seasonCode)
	}

	if bendrn != "" {
		query = query + fmt.Sprintf(" and BENDRN = '%s' ", bendrn)
		queryWithCount = queryWithCount + fmt.Sprintf(" and BENDRN = '%s' ", bendrn)
	}

	countQuery, err := s.db.Prepare(queryWithCount)
	if err != nil {
		return siteDataEntryWithCount, err
	}

	countrows, err := countQuery.Query(year, fieldOfficeCode)
	if err != nil {
		return siteDataEntryWithCount, err
	}
	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&siteDataEntryWithCount.TotalCount)
		if err != nil {
			return siteDataEntryWithCount, err
		}
	}

	siteEntries := []models.Sites{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "site_id"
	}
	fishDataEntriesSqlWithSearch := query + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(fishDataEntriesSqlWithSearch)
	if err != nil {
		return siteDataEntryWithCount, err
	}

	rows, err := dbQuery.Query(year, fieldOfficeCode)
	if err != nil {
		return siteDataEntryWithCount, err
	}
	defer rows.Close()

	for rows.Next() {
		siteDataEntry := models.Sites{}
		err = rows.Scan(&siteDataEntry.BendRiverMile, &siteDataEntry.SiteID, &siteDataEntry.SiteFID, &siteDataEntry.SiteYear, &siteDataEntry.FieldOffice, &siteDataEntry.Project,
			&siteDataEntry.Segment, &siteDataEntry.Season, &siteDataEntry.SampleUnitTypeCode, &siteDataEntry.Bendrn, &siteDataEntry.EditInitials, &siteDataEntry.UploadedBy)
		if err != nil {
			return siteDataEntryWithCount, err
		}
		siteEntries = append(siteEntries, siteDataEntry)
	}

	siteDataEntryWithCount.Items = siteEntries

	return siteDataEntryWithCount, err
}

var insertSiteDataSql = `insert into ds_sites (brm_id, site_fid, year, FIELDOFFICE, PROJECT_ID,
	SEGMENT_ID, SEASON, SAMPLE_UNIT_TYPE, BENDRN, edit_initials, last_updated, uploaded_by) values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12) returning site_id into :13`

func (s *PallidSturgeonStore) SaveSiteDataEntry(sitehDataEntry models.Sites) (int, error) {
	var id int
	_, err := s.db.Exec(insertSiteDataSql, sitehDataEntry.BendRiverMile, sitehDataEntry.SiteFID, sitehDataEntry.SiteYear, sitehDataEntry.FieldOffice, sitehDataEntry.Project,
		sitehDataEntry.Segment, sitehDataEntry.Season, sitehDataEntry.SampleUnitTypeCode, sitehDataEntry.Bendrn, sitehDataEntry.EditInitials, sitehDataEntry.LastUpdated, sitehDataEntry.UploadedBy, sql.Out{Dest: &id})

	return id, err
}

var updateSiteDataSql = `UPDATE ds_sites
SET   site_fid = :2,
	  year = :3,
	  FIELDOFFICE = :4,
	  PROJECT_ID = :5,
	  SEGMENT_ID = :6,
	  SEASON = :7,
	  SAMPLE_UNIT_TYPE = :8,
	  BENDRN = :9,
	  edit_initials = :10,
	  last_updated = :11, 
	  uploaded_by = :12,
	  brm_id = :13
WHERE site_id = :1`

func (s *PallidSturgeonStore) UpdateSiteDataEntry(sitehDataEntry models.Sites) error {
	_, err := s.db.Exec(updateSiteDataSql, sitehDataEntry.SiteFID, sitehDataEntry.SiteYear, sitehDataEntry.FieldOffice, sitehDataEntry.Project,
		sitehDataEntry.Segment, sitehDataEntry.Season, sitehDataEntry.SampleUnitTypeCode, sitehDataEntry.Bendrn, sitehDataEntry.EditInitials, sitehDataEntry.LastUpdated, sitehDataEntry.UploadedBy, sitehDataEntry.BendRiverMile, sitehDataEntry.SiteID)
	return err
}

var fishDataEntriesSql = `SELECT f_id,f_fid,FIELDOFFICE,PROJECT_ID,SEGMENT_ID,uniqueidentifier,id,panelhook,bait,SPECIES_ID,length,weight,FISHCOUNT,otolith,rayspine,scale,FTPREFIX,FTNUM,FTMR,mr_id,edit_initials,last_edit_comment, uploaded_by fROM ds_fish where FIELDOFFICE = :1`

var fishDataEntriesCountSql = `SELECT count(*) FROM ds_fish where FIELDOFFICE = :1`

var fishDataEntriesByFidSql = `SELECT f_id,f_fid,FIELDOFFICE,PROJECT_ID,SEGMENT_ID,uniqueidentifier,id,panelhook,bait,SPECIES_ID,length,weight,FISHCOUNT,otolith,rayspine,scale,FTPREFIX,FTNUM,FTMR,mr_id,edit_initials,last_edit_comment, uploaded_by fROM ds_fish where f_id = :1 and FIELDOFFICE = :2`

var fishDataEntriesCountByFidSql = `SELECT count(*) FROM ds_fish where f_id = :1 and FIELDOFFICE = :2`

var fishDataEntriesByFfidSql = `SELECT f_id,f_fid,FIELDOFFICE,PROJECT_ID,SEGMENT_ID,uniqueidentifier,id,panelhook,bait,SPECIES_ID,length,weight,FISHCOUNT,otolith,rayspine,scale,FTPREFIX,FTNUM,FTMR,mr_id,edit_initials,last_edit_comment, uploaded_by FROM ds_fish where f_fid = :1 and FIELDOFFICE = :2`

var fishDataEntriesCountByFfidSql = `SELECT count(*) FROM ds_fish where f_fid = :1 and FIELDOFFICE = :2`

var fishDataEntriesByMridSql = `SELECT f_id,f_fid,FIELDOFFICE,PROJECT_ID,SEGMENT_ID,uniqueidentifier,id,panelhook,bait,SPECIES_ID,length,weight,FISHCOUNT,otolith,rayspine,scale,FTPREFIX,FTNUM,FTMR,mr_id,edit_initials,last_edit_comment, uploaded_by FROM ds_fish where mr_id = :1 and FIELDOFFICE = :2`

var fishDataEntriesCountByMridSql = `SELECT count(*) FROM ds_fish where mr_id = :1 and FIELDOFFICE = :2`

func (s *PallidSturgeonStore) GetFishDataEntries(tableId string, fieldId string, mrId string, officeCode string, queryParams models.SearchParams) (models.FishDataEntryWithCount, error) {
	fishDataEntryWithCount := models.FishDataEntryWithCount{}
	query := ""
	queryWithCount := ""
	id := ""

	if tableId != "" {
		query = fishDataEntriesByFidSql
		queryWithCount = fishDataEntriesCountByFidSql
		id = tableId
	}

	if fieldId != "" {
		query = fishDataEntriesByFfidSql
		queryWithCount = fishDataEntriesCountByFfidSql
		id = fieldId
	}

	if mrId != "" {
		query = fishDataEntriesByMridSql
		queryWithCount = fishDataEntriesCountByMridSql
		id = mrId
	}

	if fieldId == "" && tableId == "" && mrId == "" {
		query = fishDataEntriesSql
		queryWithCount = fishDataEntriesCountSql
	}

	countQuery, err := s.db.Prepare(queryWithCount)
	if err != nil {
		return fishDataEntryWithCount, err
	}

	var countrows *sql.Rows
	if id == "" {
		countrows, err = countQuery.Query(officeCode)
		if err != nil {
			return fishDataEntryWithCount, err
		}
	} else {
		countrows, err = countQuery.Query(id, officeCode)
		if err != nil {
			return fishDataEntryWithCount, err
		}
	}

	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&fishDataEntryWithCount.TotalCount)
		if err != nil {
			return fishDataEntryWithCount, err
		}
	}

	fishEntries := []models.UploadFish{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "f_id"
	}
	fishDataEntriesSqlWithSearch := query + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(fishDataEntriesSqlWithSearch)
	if err != nil {
		return fishDataEntryWithCount, err
	}

	var rows *sql.Rows
	if id == "" {
		rows, err = dbQuery.Query(officeCode)
		if err != nil {
			return fishDataEntryWithCount, err
		}
	} else {
		rows, err = dbQuery.Query(id, officeCode)
		if err != nil {
			return fishDataEntryWithCount, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		fishDataEntry := models.UploadFish{}
		err = rows.Scan(&fishDataEntry.Fid, &fishDataEntry.Ffid, &fishDataEntry.Fieldoffice, &fishDataEntry.Project, &fishDataEntry.Segment, &fishDataEntry.UniqueID, &fishDataEntry.Id, &fishDataEntry.Panelhook,
			&fishDataEntry.Bait, &fishDataEntry.Species, &fishDataEntry.Length, &fishDataEntry.Weight, &fishDataEntry.Fishcount, &fishDataEntry.Otolith, &fishDataEntry.Rayspine,
			&fishDataEntry.Scale, &fishDataEntry.Ftprefix, &fishDataEntry.Ftnum, &fishDataEntry.Ftmr, &fishDataEntry.MrID, &fishDataEntry.EditInitials, &fishDataEntry.LastEditComment, &fishDataEntry.UploadedBy)
		if err != nil {
			return fishDataEntryWithCount, err
		}
		fishEntries = append(fishEntries, fishDataEntry)
	}

	fishDataEntryWithCount.Items = fishEntries

	return fishDataEntryWithCount, err
}

var insertFishDataSql = `insert into ds_fish (f_fid,FIELDOFFICE,PROJECT_ID,SEGMENT_ID,uniqueidentifier,id,panelhook,bait,SPECIES_ID,length,weight,FISHCOUNT,otolith,rayspine,scale,FTPREFIX,FTNUM,FTMR,mr_id,edit_initials,last_edit_comment, last_updated, uploaded_by) values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23) returning f_id into :24`

func (s *PallidSturgeonStore) SaveFishDataEntry(fishDataEntry models.UploadFish) (int, error) {
	var id int
	_, err := s.db.Exec(insertFishDataSql, fishDataEntry.Ffid, fishDataEntry.Fieldoffice, fishDataEntry.Project, fishDataEntry.Segment, fishDataEntry.UniqueID, fishDataEntry.Id, fishDataEntry.Panelhook,
		fishDataEntry.Bait, fishDataEntry.Species, fishDataEntry.Length, fishDataEntry.Weight, fishDataEntry.Fishcount, fishDataEntry.Otolith, fishDataEntry.Rayspine,
		fishDataEntry.Scale, fishDataEntry.Ftprefix, fishDataEntry.Ftnum, fishDataEntry.Ftmr, fishDataEntry.MrID, fishDataEntry.EditInitials, fishDataEntry.LastEditComment, fishDataEntry.LastUpdated, fishDataEntry.UploadedBy, sql.Out{Dest: &id})

	return id, err
}

var updateFishDataSql = `UPDATE ds_fish
SET   f_fid = :2,
	  FIELDOFFICE = :3,
	  PROJECT_ID = :4,
	  SEGMENT_ID = :5,
	  uniqueidentifier = :6,
	  id = :7,
	  panelhook = :8,
	  bait = :9,
	  SPECIES_ID = :10,
	  length = :11,
	  weight = :12,
	  FISHCOUNT = :13,
	  otolith = :14,
	  rayspine = :15,
	  scale = :16,
	  FTPREFIX = :17,
	  FTNUM = :18,
	  FTMR = :19,
	  edit_initials = :20,
	  last_edit_comment = :21,
	  last_updated = :22, 
	  uploaded_by = :23
WHERE f_id = :1`

func (s *PallidSturgeonStore) UpdateFishDataEntry(fishDataEntry models.UploadFish) error {
	_, err := s.db.Exec(updateFishDataSql, fishDataEntry.Ffid, fishDataEntry.Fieldoffice, fishDataEntry.Project, fishDataEntry.Segment, fishDataEntry.UniqueID, fishDataEntry.Id, fishDataEntry.Panelhook,
		fishDataEntry.Bait, fishDataEntry.Species, fishDataEntry.Length, fishDataEntry.Weight, fishDataEntry.Fishcount, fishDataEntry.Otolith, fishDataEntry.Rayspine,
		fishDataEntry.Scale, fishDataEntry.Ftprefix, fishDataEntry.Ftnum, fishDataEntry.Ftmr, fishDataEntry.EditInitials, fishDataEntry.LastEditComment, fishDataEntry.LastUpdated, fishDataEntry.UploadedBy, fishDataEntry.Fid)
	return err
}

var insertMoriverDataSql = `insert into ds_moriver(mr_fid,site_id,FIELDOFFICE,PROJECT_ID,SEGMENT_ID,SEASON,set_date, subsample, subsample_pass, 
	subsample_r_or_n, recorder, gear_code, GEAR_TYPE, temp, turbidity, conductivity, do,
	distance, width, net_river_mile, structure_number, usgs, river_stage, discharge,
	u1, u2, u3, u4, u5, u6, u7, MACRO_ID, MESO_ID, habitat_r_or_n, qc,
	micro_structure, structure_flow, structure_mod, set_site_1, set_site_2, set_site_3,
	start_time, start_latitude, start_longitude, stop_time, stop_latitude, stop_longitude, 
	depth_1, velocity_bottom_1, velocity_mid_1, velocity_top_1,
	depth_2, velocity_bottom_2, velocity_mid_2, velocity_top_2,
	depth_3, velocity_bottom_3, velocity_mid_3, velocity_top_3, 
	water_velocity, cobble_estimation_code, ORGANIC, silt, sand, gravel,
	comments, complete, checkby, turbidity_ind, velocity_ind, edit_initials,last_edit_comment, last_updated, uploaded_by) values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,
		:21,:22,:23,:24,:25,:26,:27,:28,:29,:30,:31,:32,:33,:34,:35,:36,:37,:38,:39,:40,:41,:42,:43,:44,:45,:46,:47,:48,:49,:50,
		:51,:52,:53,:54,:55,:56,:57,:58,:59,:60,:61,:62,:63,:64,:65,:66,:67,:68,:69,:70,:71,:72,:73,:74) returning mr_id into :75`

func (s *PallidSturgeonStore) SaveMoriverDataEntry(moriverDataEntry models.UploadMoriver) (int, error) {
	var id int
	_, err := s.db.Exec(insertMoriverDataSql, moriverDataEntry.MrFid, moriverDataEntry.SiteID, moriverDataEntry.FieldOffice,
		moriverDataEntry.Project, moriverDataEntry.Segment, moriverDataEntry.Season, moriverDataEntry.Setdate, moriverDataEntry.Subsample, moriverDataEntry.Subsamplepass,
		moriverDataEntry.SubsampleROrN, moriverDataEntry.Recorder, moriverDataEntry.Gear, moriverDataEntry.GearType, moriverDataEntry.Temp, moriverDataEntry.Turbidity, moriverDataEntry.Conductivity, moriverDataEntry.Do,
		moriverDataEntry.Distance, moriverDataEntry.Width, moriverDataEntry.Netrivermile, moriverDataEntry.Structurenumber, moriverDataEntry.Usgs, moriverDataEntry.Riverstage, moriverDataEntry.Discharge,
		moriverDataEntry.U1, moriverDataEntry.U2, moriverDataEntry.U3, moriverDataEntry.U4, moriverDataEntry.U5, moriverDataEntry.U6, moriverDataEntry.U7, moriverDataEntry.Macro, moriverDataEntry.Meso, moriverDataEntry.Habitatrn, moriverDataEntry.Qc,
		moriverDataEntry.MicroStructure, moriverDataEntry.StructureFlow, moriverDataEntry.StructureMod, moriverDataEntry.SetSite1, moriverDataEntry.SetSite2, moriverDataEntry.SetSite3,
		moriverDataEntry.StartTime, moriverDataEntry.StartLatitude, moriverDataEntry.StartLongitude, moriverDataEntry.StopTime, moriverDataEntry.StopLatitude, moriverDataEntry.StopLongitude,
		moriverDataEntry.Depth1, moriverDataEntry.Velocitybot1, moriverDataEntry.Velocity08_1, moriverDataEntry.Velocity02or06_1,
		moriverDataEntry.Depth2, moriverDataEntry.Velocitybot2, moriverDataEntry.Velocity08_2, moriverDataEntry.Velocity02or06_2,
		moriverDataEntry.Depth3, moriverDataEntry.Velocitybot3, moriverDataEntry.Velocity08_3, moriverDataEntry.Velocity02or06_3,
		moriverDataEntry.Watervel, moriverDataEntry.Cobble, moriverDataEntry.Organic, moriverDataEntry.Silt, moriverDataEntry.Sand, moriverDataEntry.Gravel,
		moriverDataEntry.Comments, moriverDataEntry.Complete, moriverDataEntry.Checkby, moriverDataEntry.NoTurbidity, moriverDataEntry.NoVelocity, moriverDataEntry.EditInitials, moriverDataEntry.LastEditComment, moriverDataEntry.LastUpdated, moriverDataEntry.UploadedBy, sql.Out{Dest: &id})
	return id, err
}

var updateMoriverDataSql = `UPDATE ds_moriver
SET  PROJECT_ID = :2,SEGMENT_ID = :3,SEASON = :4,set_date = :5, subsample = :6, subsample_pass = :7, 
	subsample_r_or_n = :8, recorder = :9, gear_code = :10, GEAR_TYPE = :11, temp = :12, turbidity = :13, conductivity = :14, do = :15,
	distance = :16, width = :17, net_river_mile = :18, structure_number = :19, usgs = :20, river_stage = :21, discharge = :22,
	u1 = :23, u2 = :24, u3 = :25, u4 = :26, u5 = :27, u6 = :28, u7 = :29, MACRO_ID = :30, MESO_ID = :31, habitat_r_or_n = :32, qc = :33,
	micro_structure = :34, structure_flow = :35, structure_mod = :36, set_site_1 = :37, set_site_2 = :38, set_site_3 = :39,
	start_time = :40, start_latitude = :41, start_longitude = :42, stop_time = :43, stop_latitude = :44, stop_longitude = :45, 
	depth_1 = :46, velocity_bottom_1 = :47, velocity_mid_1 = :48, velocity_top_1 = :49,
	depth_2 = :50, velocity_bottom_2 = :51, velocity_mid_2 = :52, velocity_top_2 = :53,
	depth_3 = :54, velocity_bottom_3 = :55, velocity_mid_3 = :56, velocity_top_3 = :57, 
	water_velocity = :58, cobble_estimation_code = :59, ORGANIC = :60, silt = :61, sand = :62, gravel = :63,
	comments = :64, complete = :65, checkby = :66, turbidity_ind = :67, velocity_ind = :68, edit_initials = :69,  mr_fid= :70, site_id = :71, FIELDOFFICE = :72,
	last_edit_comment = :73,
	last_updated = :74, 
	uploaded_by = :75
WHERE mr_id = :1`

func (s *PallidSturgeonStore) UpdateMoriverDataEntry(moriverDataEntry models.UploadMoriver) error {
	_, err := s.db.Exec(updateMoriverDataSql,
		moriverDataEntry.Project, moriverDataEntry.Segment, moriverDataEntry.Season, moriverDataEntry.Setdate, moriverDataEntry.Subsample, moriverDataEntry.Subsamplepass,
		moriverDataEntry.SubsampleROrN, moriverDataEntry.Recorder, moriverDataEntry.Gear, moriverDataEntry.GearType, moriverDataEntry.Temp, moriverDataEntry.Turbidity, moriverDataEntry.Conductivity, moriverDataEntry.Do,
		moriverDataEntry.Distance, moriverDataEntry.Width, moriverDataEntry.Netrivermile, moriverDataEntry.Structurenumber, moriverDataEntry.Usgs, moriverDataEntry.Riverstage, moriverDataEntry.Discharge,
		moriverDataEntry.U1, moriverDataEntry.U2, moriverDataEntry.U3, moriverDataEntry.U4, moriverDataEntry.U5, moriverDataEntry.U6, moriverDataEntry.U7, moriverDataEntry.Macro, moriverDataEntry.Meso, moriverDataEntry.Habitatrn, moriverDataEntry.Qc,
		moriverDataEntry.MicroStructure, moriverDataEntry.StructureFlow, moriverDataEntry.StructureMod, moriverDataEntry.SetSite1, moriverDataEntry.SetSite2, moriverDataEntry.SetSite3,
		moriverDataEntry.StartTime, moriverDataEntry.StartLatitude, moriverDataEntry.StartLongitude, moriverDataEntry.StopTime, moriverDataEntry.StopLatitude, moriverDataEntry.StopLongitude,
		moriverDataEntry.Depth1, moriverDataEntry.Velocitybot1, moriverDataEntry.Velocity08_1, moriverDataEntry.Velocity02or06_1,
		moriverDataEntry.Depth2, moriverDataEntry.Velocitybot2, moriverDataEntry.Velocity08_2, moriverDataEntry.Velocity02or06_2,
		moriverDataEntry.Depth3, moriverDataEntry.Velocitybot3, moriverDataEntry.Velocity08_3, moriverDataEntry.Velocity02or06_3,
		moriverDataEntry.Watervel, moriverDataEntry.Cobble, moriverDataEntry.Organic, moriverDataEntry.Silt, moriverDataEntry.Sand, moriverDataEntry.Gravel,
		moriverDataEntry.Comments, moriverDataEntry.Complete, moriverDataEntry.Checkby, moriverDataEntry.NoTurbidity, moriverDataEntry.NoVelocity, moriverDataEntry.EditInitials, moriverDataEntry.MrFid, moriverDataEntry.SiteID, moriverDataEntry.FieldOffice, moriverDataEntry.LastEditComment, moriverDataEntry.LastUpdated, moriverDataEntry.UploadedBy, moriverDataEntry.MrID)
	return err
}

var moriverDataEntriesByFidSql = `select mr_fid,mr_id,site_id,FIELDOFFICE,PROJECT_ID,SEGMENT_ID,SEASON,set_date, subsample, subsample_pass, 
									subsample_r_or_n, recorder, gear_code, GEAR_TYPE, temp, turbidity, conductivity, do,
									distance, width, net_river_mile, structure_number, usgs, river_stage, discharge,
									u1, u2, u3, u4, u5, u6, u7, MACRO_ID, MESO_ID, habitat_r_or_n, qc,
									micro_structure, structure_flow, structure_mod, set_site_1, set_site_2, set_site_3,
									start_time, start_latitude, start_longitude, stop_time, stop_latitude, stop_longitude, 
									depth_1, velocity_bottom_1, velocity_mid_1, velocity_top_1,
									depth_2, velocity_bottom_2, velocity_mid_2, velocity_top_2,
									depth_3, velocity_bottom_3, velocity_mid_3, velocity_top_3, 
									water_velocity, cobble_estimation_code, ORGANIC, silt, sand, gravel,
									comments, complete, checkby, turbidity_ind, velocity_ind, edit_initials,last_edit_comment, uploaded_by from ds_moriver where mr_id = :1 and FIELDOFFICE = :2`

var moriverDataEntriesCountByFidSql = `SELECT count(*) FROM ds_moriver where mr_id = :1 and FIELDOFFICE = :2`

var moriverDataEntriesByFfidSql = `select mr_fid,mr_id,site_id,FIELDOFFICE,PROJECT_ID,SEGMENT_ID,SEASON,set_date, subsample, subsample_pass, 
									subsample_r_or_n, recorder, gear_code, GEAR_TYPE, temp, turbidity, conductivity, do,
									distance, width, net_river_mile, structure_number, usgs, river_stage, discharge,
									u1, u2, u3, u4, u5, u6, u7, MACRO_ID, MESO_ID, habitat_r_or_n, qc,
									micro_structure, structure_flow, structure_mod, set_site_1, set_site_2, set_site_3,
									start_time, start_latitude, start_longitude, stop_time, stop_latitude, stop_longitude, 
									depth_1, velocity_bottom_1, velocity_mid_1, velocity_top_1,
									depth_2, velocity_bottom_2, velocity_mid_2, velocity_top_2,
									depth_3, velocity_bottom_3, velocity_mid_3, velocity_top_3, 
									water_velocity, cobble_estimation_code, ORGANIC, silt, sand, gravel,
									comments, complete, checkby, turbidity_ind, velocity_ind, edit_initials,last_edit_comment, uploaded_by from ds_moriver where mr_fid = :1 and FIELDOFFICE = :2`

var moriverDataEntriesCountByFfidSql = `SELECT count(*) FROM ds_moriver where mr_fid = :1 and FIELDOFFICE = :2`

func (s *PallidSturgeonStore) GetMoriverDataEntries(tableId string, fieldId string, fieldOfficeCode string, queryParams models.SearchParams) (models.MoriverDataEntryWithCount, error) {
	moriverDataEntryWithCount := models.MoriverDataEntryWithCount{}
	query := ""
	queryWithCount := ""
	id := ""

	if tableId != "" {
		query = moriverDataEntriesByFidSql
		queryWithCount = moriverDataEntriesCountByFidSql
		id = tableId
	}

	if fieldId != "" {
		query = moriverDataEntriesByFfidSql
		queryWithCount = moriverDataEntriesCountByFfidSql
		id = fieldId
	}

	countQuery, err := s.db.Prepare(queryWithCount)
	if err != nil {
		return moriverDataEntryWithCount, err
	}

	countrows, err := countQuery.Query(id, fieldOfficeCode)
	if err != nil {
		return moriverDataEntryWithCount, err
	}
	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&moriverDataEntryWithCount.TotalCount)
		if err != nil {
			return moriverDataEntryWithCount, err
		}
	}

	moriverEntries := []models.UploadMoriver{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "mr_id"
	}
	moriverDataEntriesSqlWithSearch := query + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(moriverDataEntriesSqlWithSearch)
	if err != nil {
		return moriverDataEntryWithCount, err
	}

	rows, err := dbQuery.Query(id, fieldOfficeCode)
	if err != nil {
		return moriverDataEntryWithCount, err
	}
	defer rows.Close()

	for rows.Next() {
		moriverDataEntry := models.UploadMoriver{}
		err = rows.Scan(&moriverDataEntry.MrFid, &moriverDataEntry.MrID, &moriverDataEntry.SiteID, &moriverDataEntry.FieldOffice,
			&moriverDataEntry.Project, &moriverDataEntry.Segment, &moriverDataEntry.Season, &moriverDataEntry.Setdate, &moriverDataEntry.Subsample, &moriverDataEntry.Subsamplepass,
			&moriverDataEntry.SubsampleROrN, &moriverDataEntry.Recorder, &moriverDataEntry.Gear, &moriverDataEntry.GearType, &moriverDataEntry.Temp, &moriverDataEntry.Turbidity, &moriverDataEntry.Conductivity, &moriverDataEntry.Do,
			&moriverDataEntry.Distance, &moriverDataEntry.Width, &moriverDataEntry.Netrivermile, &moriverDataEntry.Structurenumber, &moriverDataEntry.Usgs, &moriverDataEntry.Riverstage, &moriverDataEntry.Discharge,
			&moriverDataEntry.U1, &moriverDataEntry.U2, &moriverDataEntry.U3, &moriverDataEntry.U4, &moriverDataEntry.U5, &moriverDataEntry.U6, &moriverDataEntry.U7, &moriverDataEntry.Macro, &moriverDataEntry.Meso, &moriverDataEntry.Habitatrn, &moriverDataEntry.Qc,
			&moriverDataEntry.MicroStructure, &moriverDataEntry.StructureFlow, &moriverDataEntry.StructureMod, &moriverDataEntry.SetSite1, &moriverDataEntry.SetSite2, &moriverDataEntry.SetSite3,
			&moriverDataEntry.StartTime, &moriverDataEntry.StartLatitude, &moriverDataEntry.StartLongitude, &moriverDataEntry.StopTime, &moriverDataEntry.StopLatitude, &moriverDataEntry.StopLongitude,
			&moriverDataEntry.Depth1, &moriverDataEntry.Velocitybot1, &moriverDataEntry.Velocity08_1, &moriverDataEntry.Velocity02or06_1,
			&moriverDataEntry.Depth2, &moriverDataEntry.Velocitybot2, &moriverDataEntry.Velocity08_2, &moriverDataEntry.Velocity02or06_2,
			&moriverDataEntry.Depth3, &moriverDataEntry.Velocitybot3, &moriverDataEntry.Velocity08_3, &moriverDataEntry.Velocity02or06_3,
			&moriverDataEntry.Watervel, &moriverDataEntry.Cobble, &moriverDataEntry.Organic, &moriverDataEntry.Silt, &moriverDataEntry.Sand, &moriverDataEntry.Gravel,
			&moriverDataEntry.Comments, &moriverDataEntry.Complete, &moriverDataEntry.Checkby, &moriverDataEntry.NoTurbidity, &moriverDataEntry.NoVelocity, &moriverDataEntry.EditInitials, &moriverDataEntry.LastEditComment, &moriverDataEntry.UploadedBy)
		if err != nil {
			return moriverDataEntryWithCount, err
		}
		moriverEntries = append(moriverEntries, moriverDataEntry)
	}

	moriverDataEntryWithCount.Items = moriverEntries

	return moriverDataEntryWithCount, err
}

var insertSupplementalDataSql = `insert into ds_supplemental(f_id, f_fid, mr_id,
	TAGNUMBER, PITRN,
	SCUTELOC, SCUTENUM, SCUTELOC2, SCUTENUM_2,
	ELHV, ELCOLOR, ERHV, ERCOLOR, CWTYN, DANGLER, genetic_y_n_or_u, genetics_vial_number,
	BROODSTOCK, HATCH_WILD, species_id,
	head, snouttomouth, inter, mouthwidth, m_ib,
	l_ob, l_ib, r_ib,
	r_ob, anal, dorsal, status, HATCHERY_ORIGIN,
	SEX, stage,  recapture, photo,
	genetic_needs, other_tag_info,
	comments,edit_initials,last_edit_comment, last_updated, uploaded_by) values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,
		:21,:22,:23,:24,:25,:26,:27,:28,:29,:30,:31,:32,:33,:34,:35,:36,:37,:38,:39,:40,:41,:42,:43,:44, :45) returning s_id into :46`

func (s *PallidSturgeonStore) SaveSupplementalDataEntry(supplementalDataEntry models.UploadSupplemental) (int, error) {
	var id int
	_, err := s.db.Exec(insertSupplementalDataSql,
		supplementalDataEntry.Fid,
		supplementalDataEntry.FFid,
		supplementalDataEntry.MrId,
		supplementalDataEntry.Tagnumber,
		supplementalDataEntry.Pitrn,
		supplementalDataEntry.Scuteloc,
		supplementalDataEntry.Scutenum,
		supplementalDataEntry.Scuteloc2,
		supplementalDataEntry.Scutenum2,
		supplementalDataEntry.Elhv,
		supplementalDataEntry.Elcolor,
		supplementalDataEntry.Erhv,
		supplementalDataEntry.Ercolor,
		supplementalDataEntry.Cwtyn,
		supplementalDataEntry.Dangler,
		supplementalDataEntry.Genetic,
		supplementalDataEntry.GeneticsVialNumber,
		supplementalDataEntry.Broodstock,
		supplementalDataEntry.HatchWild,
		supplementalDataEntry.SpeciesId,
		supplementalDataEntry.Head,
		supplementalDataEntry.Snouttomouth,
		supplementalDataEntry.Inter,
		supplementalDataEntry.Mouthwidth,
		supplementalDataEntry.MIb,
		supplementalDataEntry.LOb,
		supplementalDataEntry.LIb,
		supplementalDataEntry.RIb,
		supplementalDataEntry.ROb,
		supplementalDataEntry.Anal,
		supplementalDataEntry.Dorsal,
		supplementalDataEntry.Status,
		supplementalDataEntry.HatcheryOrigin,
		supplementalDataEntry.Sex,
		supplementalDataEntry.Stage,
		supplementalDataEntry.Recapture,
		supplementalDataEntry.Photo,
		supplementalDataEntry.GeneticNeeds,
		supplementalDataEntry.OtherTagInfo,
		supplementalDataEntry.Comments,
		supplementalDataEntry.EditInitials,
		supplementalDataEntry.LastEditComment,
		supplementalDataEntry.LastUpdated,
		supplementalDataEntry.UploadedBy,
		sql.Out{Dest: &id})
	return id, err
}

var updateSupplementalDataSql = `UPDATE ds_supplemental
SET   f_fid = :2,mr_id = :3,TAGNUMBER = :4, PITRN = :5, SCUTELOC = :6, 
		SCUTENUM = :7, SCUTELOC2  = :8, SCUTENUM_2 = :9, ELHV = :10, ELCOLOR = :11, ERHV = :12, 
		ERCOLOR = :13, CWTYN  = :14, DANGLER = :15, genetic_y_n_or_u = :16, genetics_vial_number = :17,
		BROODSTOCK = :18, HATCH_WILD = :19, 
		head = :20, snouttomouth = :21, inter = :22, mouthwidth = :23, m_ib = :24, l_ob = :25, l_ib = :26, r_ib = :27, 
		r_ob = :28, anal = :29, dorsal = :30, status = :31, HATCHERY_ORIGIN = :32, SEX = :33, stage = :34,  recapture = :35, 
		photo = :36, genetic_needs = :37, other_tag_info = :38, comments = :39, f_id = :40,
		edit_initials = :41,
		last_edit_comment = :42,	
		last_updated = :43, 
		uploaded_by = :44,
		species_id = :45
WHERE f_id = :1`

func (s *PallidSturgeonStore) UpdateSupplementalDataEntry(supplementalDataEntry models.UploadSupplemental) error {
	_, err := s.db.Exec(updateSupplementalDataSql,
		supplementalDataEntry.FFid,
		supplementalDataEntry.MrId,
		supplementalDataEntry.Tagnumber,
		supplementalDataEntry.Pitrn,
		supplementalDataEntry.Scuteloc,
		supplementalDataEntry.Scutenum,
		supplementalDataEntry.Scuteloc2,
		supplementalDataEntry.Scutenum2,
		supplementalDataEntry.Elhv,
		supplementalDataEntry.Elcolor,
		supplementalDataEntry.Erhv,
		supplementalDataEntry.Ercolor,
		supplementalDataEntry.Cwtyn,
		supplementalDataEntry.Dangler,
		supplementalDataEntry.Genetic,
		supplementalDataEntry.GeneticsVialNumber,
		supplementalDataEntry.Broodstock,
		supplementalDataEntry.HatchWild,
		supplementalDataEntry.Head,
		supplementalDataEntry.Snouttomouth,
		supplementalDataEntry.Inter,
		supplementalDataEntry.Mouthwidth,
		supplementalDataEntry.MIb,
		supplementalDataEntry.LOb,
		supplementalDataEntry.LIb,
		supplementalDataEntry.RIb,
		supplementalDataEntry.ROb,
		supplementalDataEntry.Anal,
		supplementalDataEntry.Dorsal,
		supplementalDataEntry.Status,
		supplementalDataEntry.HatcheryOrigin,
		supplementalDataEntry.Sex,
		supplementalDataEntry.Stage,
		supplementalDataEntry.Recapture,
		supplementalDataEntry.Photo,
		supplementalDataEntry.GeneticNeeds,
		supplementalDataEntry.OtherTagInfo,
		supplementalDataEntry.Comments,
		supplementalDataEntry.Fid,
		supplementalDataEntry.EditInitials,
		supplementalDataEntry.LastEditComment,
		supplementalDataEntry.LastUpdated,
		supplementalDataEntry.UploadedBy,
		supplementalDataEntry.Sid,
		supplementalDataEntry.SpeciesId)
	return err
}

var supplementalDataEntriesByFidSql = `select f_id, f_fid, mr_id,
										TAGNUMBER, PITRN, 
										SCUTELOC, SCUTENUM, SCUTELOC2, SCUTENUM_2, 
										ELHV, ELCOLOR, ERHV, ERCOLOR, CWTYN, DANGLER, genetic_y_n_or_u, genetics_vial_number,
										BROODSTOCK, HATCH_WILD, species_id,
										head, snouttomouth, inter, mouthwidth, m_ib,
										l_ob, l_ib, r_ib, 
										r_ob, anal, dorsal, status, HATCHERY_ORIGIN, 
										SEX, stage,  recapture, photo,
										genetic_needs, other_tag_info,
										comments, edit_initials,last_edit_comment, uploaded_by from ds_supplemental where f_id = :1 `

var supplementalDataEntriesCountByFidSql = `SELECT count(*) FROM ds_supplemental where f_id = :1`

var supplementalDataEntriesByFfidSql = `select f_id, f_fid, mr_id,
										TAGNUMBER, PITRN, 
										SCUTELOC, SCUTENUM, SCUTELOC2, SCUTENUM_2, 
										ELHV, ELCOLOR, ERHV, ERCOLOR, CWTYN, DANGLER, genetic_y_n_or_u, genetics_vial_number,
										BROODSTOCK, HATCH_WILD, species_id,
										head, snouttomouth, inter, mouthwidth, m_ib,
										l_ob, l_ib, r_ib, 
										r_ob, anal, dorsal, status, HATCHERY_ORIGIN, 
										SEX, stage,  recapture, photo,
										genetic_needs, other_tag_info,
										comments, edit_initials,last_edit_comment, uploaded_by from ds_supplemental where f_fid = :1`

var supplementalDataEntriesCountByFfidSql = `SELECT count(*) FROM ds_supplemental where f_fid = :1`

var supplementalDataEntriesByGeneticsVialSql = `select f_id, f_fid, mr_id,
										TAGNUMBER, PITRN, 
										SCUTELOC, SCUTENUM, SCUTELOC2, SCUTENUM_2, 
										ELHV, ELCOLOR, ERHV, ERCOLOR, CWTYN, DANGLER, genetic_y_n_or_u, genetics_vial_number,
										BROODSTOCK, HATCH_WILD, species_id,
										head, snouttomouth, inter, mouthwidth, m_ib,
										l_ob, l_ib, r_ib, 
										r_ob, anal, dorsal, status, HATCHERY_ORIGIN, 
										SEX, stage,  recapture, photo,
										genetic_needs, other_tag_info,
										comments, edit_initials,last_edit_comment, uploaded_by from ds_supplemental where genetics_vial_number = :1 `

var supplementalDataEntriesCountByGeneticsVialSql = `SELECT count(*) FROM ds_supplemental where genetics_vial_number = :1`

var supplementalDataEntriesByGeneticsPitTagSql = `select f_id, f_fid, mr_id,
										TAGNUMBER, PITRN, 
										SCUTELOC, SCUTENUM, SCUTELOC2, SCUTENUM_2, 
										ELHV, ELCOLOR, ERHV, ERCOLOR, CWTYN, DANGLER, genetic_y_n_or_u, genetics_vial_number,
										BROODSTOCK, HATCH_WILD, species_id,
										head, snouttomouth, inter, mouthwidth, m_ib,
										l_ob, l_ib, r_ib, 
										r_ob, anal, dorsal, status, HATCHERY_ORIGIN, 
										SEX, stage,  recapture, photo,
										genetic_needs, other_tag_info,
										comments, edit_initials,last_edit_comment, uploaded_by from ds_supplemental where TAGNUMBER = :1 `

var supplementalDataEntriesCountByPitTagSql = `SELECT count(*) FROM ds_supplemental where TAGNUMBER = :1 `

var supplementalDataEntriesByMrIdSql = `select f_id, f_fid, mr_id,
										TAGNUMBER, PITRN, 
										SCUTELOC, SCUTENUM, SCUTELOC2, SCUTENUM_2, 
										ELHV, ELCOLOR, ERHV, ERCOLOR, CWTYN, DANGLER, genetic_y_n_or_u, genetics_vial_number,
										BROODSTOCK, HATCH_WILD, species_id,
										head, snouttomouth, inter, mouthwidth, m_ib,
										l_ob, l_ib, r_ib, 
										r_ob, anal, dorsal, status, HATCHERY_ORIGIN, 
										SEX, stage,  recapture, photo,
										genetic_needs, other_tag_info,
										comments, edit_initials,last_edit_comment, uploaded_by from ds_supplemental where mr_id = :1 `

var supplementalDataEntriesCountByMrIdSql = `SELECT count(*) FROM ds_supplemental where mr_id = :1 `

func (s *PallidSturgeonStore) GetSupplementalDataEntries(tableId string, fieldId string, geneticsVial string, pitTag string, mrId string, queryParams models.SearchParams) (models.SupplementalDataEntryWithCount, error) {
	supplementalDataEntryWithCount := models.SupplementalDataEntryWithCount{}
	query := ""
	queryWithCount := ""
	id := ""

	if tableId != "" {
		query = supplementalDataEntriesByFidSql
		queryWithCount = supplementalDataEntriesCountByFidSql
		id = tableId
	}

	if fieldId != "" {
		query = supplementalDataEntriesByFfidSql
		queryWithCount = supplementalDataEntriesCountByFfidSql
		id = fieldId
	}

	if geneticsVial != "" {
		query = supplementalDataEntriesByGeneticsVialSql
		queryWithCount = supplementalDataEntriesCountByGeneticsVialSql
		id = geneticsVial
	}

	if pitTag != "" {
		query = supplementalDataEntriesByGeneticsPitTagSql
		queryWithCount = supplementalDataEntriesCountByPitTagSql
		id = pitTag
	}

	if mrId != "" {
		query = supplementalDataEntriesByMrIdSql
		queryWithCount = supplementalDataEntriesCountByMrIdSql
		id = mrId
	}

	countQuery, err := s.db.Prepare(queryWithCount)
	if err != nil {
		return supplementalDataEntryWithCount, err
	}

	countrows, err := countQuery.Query(id)
	if err != nil {
		return supplementalDataEntryWithCount, err
	}
	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&supplementalDataEntryWithCount.TotalCount)
		if err != nil {
			return supplementalDataEntryWithCount, err
		}
	}

	supplementalEntries := []models.UploadSupplemental{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "mr_id"
	}
	supplementalDataEntriesSqlWithSearch := query + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(supplementalDataEntriesSqlWithSearch)
	if err != nil {
		return supplementalDataEntryWithCount, err
	}

	rows, err := dbQuery.Query(id)
	if err != nil {
		return supplementalDataEntryWithCount, err
	}
	defer rows.Close()

	for rows.Next() {
		supplementalDataEntry := models.UploadSupplemental{}
		err = rows.Scan(
			&supplementalDataEntry.Fid,
			&supplementalDataEntry.FFid,
			&supplementalDataEntry.MrId,
			&supplementalDataEntry.Tagnumber,
			&supplementalDataEntry.Pitrn,
			&supplementalDataEntry.Scuteloc,
			&supplementalDataEntry.Scutenum,
			&supplementalDataEntry.Scuteloc2,
			&supplementalDataEntry.Scutenum2,
			&supplementalDataEntry.Elhv,
			&supplementalDataEntry.Elcolor,
			&supplementalDataEntry.Erhv,
			&supplementalDataEntry.Ercolor,
			&supplementalDataEntry.Cwtyn,
			&supplementalDataEntry.Dangler,
			&supplementalDataEntry.Genetic,
			&supplementalDataEntry.GeneticsVialNumber,
			&supplementalDataEntry.Broodstock,
			&supplementalDataEntry.HatchWild,
			&supplementalDataEntry.SpeciesId,
			&supplementalDataEntry.Head,
			&supplementalDataEntry.Snouttomouth,
			&supplementalDataEntry.Inter,
			&supplementalDataEntry.Mouthwidth,
			&supplementalDataEntry.MIb,
			&supplementalDataEntry.LOb,
			&supplementalDataEntry.LIb,
			&supplementalDataEntry.RIb,
			&supplementalDataEntry.ROb,
			&supplementalDataEntry.Anal,
			&supplementalDataEntry.Dorsal,
			&supplementalDataEntry.Status,
			&supplementalDataEntry.HatcheryOrigin,
			&supplementalDataEntry.Sex,
			&supplementalDataEntry.Stage,
			&supplementalDataEntry.Recapture,
			&supplementalDataEntry.Photo,
			&supplementalDataEntry.GeneticNeeds,
			&supplementalDataEntry.OtherTagInfo,
			&supplementalDataEntry.Comments,
			&supplementalDataEntry.EditInitials,
			&supplementalDataEntry.LastEditComment,
			&supplementalDataEntry.UploadedBy)
		if err != nil {
			return supplementalDataEntryWithCount, err
		}
		supplementalEntries = append(supplementalEntries, supplementalDataEntry)
	}

	supplementalDataEntryWithCount.Items = supplementalEntries

	return supplementalDataEntryWithCount, err
}

var fishDataSummaryFullDataSql = `select * FROM table (pallid_data_api.fish_datasummary_fnc(:1, :2, :3, :4, :5, :6, :7, to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

func (s *PallidSturgeonStore) GetFullFishDataSummary(year string, officeCode string, project string, approved string, season string, spice string, month string, fromDate string, toDate string) (string, error) {
	dbQuery, err := s.db.Prepare(fishDataSummaryFullDataSql)
	if err != nil {
		log.Fatal("Cannot create to file", err)
	}

	rows, err := dbQuery.Query(year, officeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		log.Fatal("Cannot create to file", err)
	}
	defer rows.Close()

	cols, _ := rows.Columns()

	file, err := os.Create("FishDataSummary.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	//save header
	data := make([]string, 0)
	data = append(data, cols...)
	err = writer.Write(data)
	if err != nil {
		log.Fatal("Cannot write to file", err)
	}

	for rows.Next() {

		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		data := make([]string, 0)
		for i := range cols {
			var v string
			val := columns[i]

			if val == nil {
				v = ""
			} else {
				v = fmt.Sprintf("%v", val)
			}
			data = append(data, v)
		}

		err := writer.Write(data)
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}

	return file.Name(), err
}

var fishDataSummarySql = `SELECT mr_id, f_id, year, FIELD_OFFICE_CODE, PROJECT_CODE, SEGMENT_CODE, SEASON_CODE, BEND_NUMBER, BEND_R_OR_N, bend_river_mile, panelhook, SPECIES_CODE, HATCHERY_ORIGIN_CODE, checkby FROM table (pallid_data_api.fish_datasummary_fnc(:1, :2, :3, :4, :5, :6, :7, to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

var fishDataSummaryCountSql = `SELECT count(*) FROM table (pallid_data_api.fish_datasummary_fnc(:1, :2, :3, :4, :5, :6, :7, to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

func (s *PallidSturgeonStore) GetFishDataSummary(year string, officeCode string, project string, approved string, season string, spice string, month string, fromDate string, toDate string, queryParams models.SearchParams) (models.FishSummaryWithCount, error) {
	fishSummariesWithCount := models.FishSummaryWithCount{}
	countQuery, err := s.db.Prepare(fishDataSummaryCountSql)
	if err != nil {
		return fishSummariesWithCount, err
	}

	countrows, err := countQuery.Query(year, officeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return fishSummariesWithCount, err
	}
	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&fishSummariesWithCount.TotalCount)
		if err != nil {
			return fishSummariesWithCount, err
		}
	}

	fishSummaries := []models.FishSummary{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "mr_id"
	}
	fishDataSummarySqlWithSearch := fishDataSummarySql + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(fishDataSummarySqlWithSearch)
	if err != nil {
		return fishSummariesWithCount, err
	}

	rows, err := dbQuery.Query(year, officeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return fishSummariesWithCount, err
	}
	defer rows.Close()

	for rows.Next() {
		fishSummary := models.FishSummary{}
		err = rows.Scan(&fishSummary.UniqueID, &fishSummary.FishID, &fishSummary.Year, &fishSummary.FieldOffice, &fishSummary.Project,
			&fishSummary.Segment, &fishSummary.Season, &fishSummary.Bend, &fishSummary.Bendrn, &fishSummary.BendRiverMile, &fishSummary.Panelhook,
			&fishSummary.Species, &fishSummary.HatcheryOrigin, &fishSummary.CheckedBy)
		if err != nil {
			return fishSummariesWithCount, err
		}
		fishSummaries = append(fishSummaries, fishSummary)
	}

	fishSummariesWithCount.Items = fishSummaries

	return fishSummariesWithCount, err
}

var suppDataSummaryFullDataSql = `select * FROM table (pallid_data_api.supp_datasummary_fnc(:1, :2, :3, :4, :5, :6, :7, to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

func (s *PallidSturgeonStore) GetFullSuppDataSummary(year string, officeCode string, project string, approved string, season string, spice string, month string, fromDate string, toDate string) (string, error) {
	dbQuery, err := s.db.Prepare(suppDataSummaryFullDataSql)
	if err != nil {
		return "Cannot create file", err
	}

	rows, err := dbQuery.Query(year, officeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return "Cannot create file", err
	}
	defer rows.Close()

	cols, _ := rows.Columns()

	file, err := os.Create("SupplementalDataSummary.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	//save header
	data := make([]string, 0)
	data = append(data, cols...)
	err = writer.Write(data)
	if err != nil {
		log.Fatal("Cannot write to file", err)
	}

	for rows.Next() {

		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		data := make([]string, 0)

		for i := range cols {
			var v string
			val := columns[i]

			if val == nil {
				v = ""
			} else {
				v = fmt.Sprintf("%v", val)
			}
			data = append(data, v)
		}

		err := writer.Write(data)
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}

	return file.Name(), err
}

var suppDataSummarySql = `SELECT fish_code, 
	mr_id, 
	f_id, 
	sid_display, 
	year, 
	FIELD_OFFICE_CODE, 
	PROJECT_CODE, 
	SEGMENT_CODE, 
	SEASON_CODE, 
	COALESCE(BEND_NUMBER, 0) as BEND_NUMBER,
	COALESCE(BEND_R_OR_N, '') as BEND_R_OR_N,
	COALESCE(bend_river_mile, 0) as bend_river_mile,
	COALESCE(HATCHERY_ORIGIN_CODE, '') as HATCHERY_ORIGIN_CODE,
	COALESCE(tag_number, '') as tag_number,
	checkby FROM table (pallid_data_api.supp_datasummary_fnc(:1, :2, :3, :4, :5, :6, :7, to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

var suppDataSummaryCountSql = `SELECT count(*) FROM table (pallid_data_api.supp_datasummary_fnc(:1, :2, :3, :4, :5, :6, :7, to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

func (s *PallidSturgeonStore) GetSuppDataSummary(year string, officeCode string, project string, approved string, season string, spice string, month string, fromDate string, toDate string, queryParams models.SearchParams) (models.SuppSummaryWithCount, error) {
	suppSummariesWithCount := models.SuppSummaryWithCount{}
	countQuery, err := s.db.Prepare(suppDataSummaryCountSql)
	if err != nil {
		return suppSummariesWithCount, err
	}

	countrows, err := countQuery.Query(year, officeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return suppSummariesWithCount, err
	}
	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&suppSummariesWithCount.TotalCount)
		if err != nil {
			return suppSummariesWithCount, err
		}
	}

	suppSummaries := []models.SuppSummary{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "mr_id"
	}
	suppDataSummarySqlWithSearch := suppDataSummarySql + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(suppDataSummarySqlWithSearch)
	if err != nil {
		return suppSummariesWithCount, err
	}

	rows, err := dbQuery.Query(year, officeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return suppSummariesWithCount, err
	}
	defer rows.Close()

	for rows.Next() {
		summary := models.SuppSummary{}
		err = rows.Scan(&summary.FishCode, &summary.UniqueID, &summary.FishID, &summary.SuppID, &summary.Year,
			&summary.FieldOffice, &summary.Project, &summary.Segment, &summary.Season, &summary.Bend, &summary.Bendrn,
			&summary.BendRiverMile, &summary.HatcheryOrigin, &summary.TagNumber, &summary.CheckedBy)
		if err != nil {
			return suppSummariesWithCount, err
		}
		suppSummaries = append(suppSummaries, summary)
	}

	suppSummariesWithCount.Items = suppSummaries

	return suppSummariesWithCount, err
}

var missouriDataSummaryFullDataSql = `SELECT * FROM table (pallid_data_api.missouri_datasummary_fnc(:1, :2, :3, :4, :5, :6, :7, to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

func (s *PallidSturgeonStore) GetFullMissouriDataSummary(year string, officeCode string, project string, approved string, season string, spice string, month string, fromDate string, toDate string) (string, error) {
	dbQuery, err := s.db.Prepare(missouriDataSummaryFullDataSql)
	if err != nil {
		return "Cannot create file", err
	}

	rows, err := dbQuery.Query(year, officeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return "Cannot create file", err
	}
	defer rows.Close()

	cols, _ := rows.Columns()

	file, err := os.Create("MissouriDataSummary.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	//save header
	data := make([]string, 0)
	data = append(data, cols...)
	err = writer.Write(data)
	if err != nil {
		log.Fatal("Cannot write to file", err)
	}

	for rows.Next() {

		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		data := make([]string, 0)
		for i := range cols {
			var v string
			val := columns[i]

			if val == nil {
				v = ""
			} else {
				v = fmt.Sprintf("%v", val)
			}
			data = append(data, v)
		}

		err := writer.Write(data)
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}

	return file.Name(), err
}

var missouriDataSummarySql = `SELECT mr_id, year, FIELD_OFFICE_CODE, PROJECT_CODE, SEGMENT_CODE, SEASON_CODE, BEND_NUMBER, BEND_R_OR_N, bend_river_mile, subsample, subsample_pass, set_Date, conductivity, checkby FROM table (pallid_data_api.missouri_datasummary_fnc(:1, :2, :3, :4, :5, :6, :7, to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

var missouriDataSummaryCountSql = `SELECT count(*) FROM table (pallid_data_api.missouri_datasummary_fnc(:1, :2, :3, :4, :5, :6, :7, to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

func (s *PallidSturgeonStore) GetMissouriDataSummary(year string, officeCode string, project string, approved string, season string, spice string, month string, fromDate string, toDate string, queryParams models.SearchParams) (models.MissouriSummaryWithCount, error) {
	missouriSummariesWithCount := models.MissouriSummaryWithCount{}
	countQuery, err := s.db.Prepare(missouriDataSummaryCountSql)
	if err != nil {
		return missouriSummariesWithCount, err
	}

	countrows, err := countQuery.Query(year, officeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return missouriSummariesWithCount, err
	}
	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&missouriSummariesWithCount.TotalCount)
		if err != nil {
			return missouriSummariesWithCount, err
		}
	}

	missouriSummaries := []models.MissouriSummary{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "mr_id"
	}
	missouriDataSummarySqlWithSearch := missouriDataSummarySql + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(missouriDataSummarySqlWithSearch)
	if err != nil {
		return missouriSummariesWithCount, err
	}

	rows, err := dbQuery.Query(year, officeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return missouriSummariesWithCount, err
	}
	defer rows.Close()

	for rows.Next() {
		summary := models.MissouriSummary{}
		err = rows.Scan(&summary.UniqueID, &summary.Year, &summary.FieldOffice, &summary.Project, &summary.Segment,
			&summary.Season, &summary.Bend, &summary.Bendrn, &summary.BendRiverMile,
			&summary.Subsample, &summary.Pass, &summary.SetDate, &summary.Conductivity, &summary.CheckedBy)
		if err != nil {
			return missouriSummariesWithCount, err
		}
		missouriSummaries = append(missouriSummaries, summary)
	}

	missouriSummariesWithCount.Items = missouriSummaries

	return missouriSummariesWithCount, err
}

var geneticDataSummaryFullDataSql = `SELECT * FROM table (pallid_data_api.genetic_datasummary_fnc(:1, :2, :3, to_date(:4,'MM/DD/YYYY'), to_date(:5,'MM/DD/YYYY'), :6, :7, :8, :9))`

func (s *PallidSturgeonStore) GetFullGeneticDataSummary(year string, officeCode string, project string, fromDate string, toDate string, broodstock string, hatchwild string, speciesid string, archive string) (string, error) {
	dbQuery, err := s.db.Prepare(geneticDataSummaryFullDataSql)
	if err != nil {
		return "Cannot create file", err
	}

	rows, err := dbQuery.Query(year, officeCode, project, fromDate, toDate, broodstock, hatchwild, hatchwild, speciesid)
	if err != nil {
		return "Cannot create file", err
	}
	defer rows.Close()

	cols, _ := rows.Columns()

	file, err := os.Create("GeneticDataSummary.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	//save header
	data := make([]string, 0)
	data = append(data, cols...)
	err = writer.Write(data)
	if err != nil {
		log.Fatal("Cannot write to file", err)
	}

	for rows.Next() {

		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		data := make([]string, 0)
		for i := range cols {
			var v string
			val := columns[i]

			if val == nil {
				v = ""
			} else {
				v = fmt.Sprintf("%v", val)
			}
			data = append(data, v)
		}

		err := writer.Write(data)
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}

	return file.Name(), err
}

var geneticDataSummarySql = `SELECT year,FIELD_OFFICE_CODE,PROJECT_CODE,genetics_vial_number,pit_tag,river,river_mile,
							state,set_date,broodstock_yn,hatchwild_yn,Speciesid_yn,archive_yn FROM table (pallid_data_api.genetic_datasummary_fnc(:1, :2, :3, to_date(:4,'MM/DD/YYYY'), to_date(:5,'MM/DD/YYYY'), :6, :7, :8, :9))`

var geneticDataSummaryCountSql = `SELECT count(*) FROM table (pallid_data_api.genetic_datasummary_fnc(:1, :2, :3, to_date(:4,'MM/DD/YYYY'), to_date(:5,'MM/DD/YYYY'), :6, :7, :8, :9))`

func (s *PallidSturgeonStore) GetGeneticDataSummary(year string, officeCode string, project string, fromDate string, toDate string, broodstock string, hatchwild string, speciesid string, archive string, queryParams models.SearchParams) (models.GeneticSummaryWithCount, error) {
	geneticSummariesWithCount := models.GeneticSummaryWithCount{}
	countQuery, err := s.db.Prepare(geneticDataSummaryCountSql)
	if err != nil {
		return geneticSummariesWithCount, err
	}

	countrows, err := countQuery.Query(year, officeCode, project, fromDate, toDate, broodstock, hatchwild, hatchwild, speciesid)
	if err != nil {
		return geneticSummariesWithCount, err
	}
	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&geneticSummariesWithCount.TotalCount)
		if err != nil {
			return geneticSummariesWithCount, err
		}
	}

	geneticSummaries := []models.GeneticSummary{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "genetics_vial_number"
	}
	geneticDataSummarySqlWithSearch := geneticDataSummarySql + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(geneticDataSummarySqlWithSearch)
	if err != nil {
		return geneticSummariesWithCount, err
	}

	rows, err := dbQuery.Query(year, officeCode, project, fromDate, toDate, broodstock, hatchwild, hatchwild, speciesid)
	if err != nil {
		return geneticSummariesWithCount, err
	}
	defer rows.Close()

	for rows.Next() {
		summary := models.GeneticSummary{}
		err = rows.Scan(&summary.Year, &summary.FieldOffice, &summary.Project, &summary.GeneticsVialNumber,
			&summary.PitTag, &summary.River, &summary.RiverMile, &summary.State, &summary.SetDate, &summary.Broodstock,
			&summary.HatchWild, &summary.SpeciesID, &summary.Archive)
		if err != nil {
			return geneticSummariesWithCount, err
		}
		geneticSummaries = append(geneticSummaries, summary)
	}

	geneticSummariesWithCount.Items = geneticSummaries

	return geneticSummariesWithCount, err
}

func (s *PallidSturgeonStore) GetFullSearchDataSummary() (string, error) {
	rows, err := s.db.Queryx("SELECT * FROM ds_search")
	if err != nil {
		return "Cannot create file", err
	}
	defer rows.Close()

	cols, _ := rows.Columns()

	file, err := os.Create("SearchDataSummary.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	//save header
	data := make([]string, 0)
	data = append(data, cols...)
	err = writer.Write(data)
	if err != nil {
		log.Fatal("Cannot write to file", err)
	}

	for rows.Next() {

		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		data := make([]string, 0)

		for i := range cols {
			var v string
			val := columns[i]

			if val == nil {
				v = ""
			} else {
				v = fmt.Sprintf("%v", val)
			}
			data = append(data, v)
		}

		err := writer.Write(data)
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}

	return file.Name(), err
}

var searchDataSummarySql = `SELECT se_id,search_date,recorder,search_type_code,start_time,start_latitude,start_longitude,stop_time,stop_latitude, stop_longitude,se_fid,ds_id,site_fid,temp,conductivity FROM ds_search`

var searchDataSummaryCountSql = `SELECT count(*) FROM ds_search`

func (s *PallidSturgeonStore) GetSearchDataSummary(queryParams models.SearchParams) (models.SearchSummaryWithCount, error) {
	searchSummariesWithCount := models.SearchSummaryWithCount{}

	filterQuery := ""
	if queryParams.Filter != "" {
		filter := "'%" + strings.ToUpper(queryParams.Filter) + "%'"
		filterQuery = fmt.Sprintf(" where se_id like %s or TO_CHAR(search_date, 'MM/DD/YYYY') like %s or UPPER(recorder) like %s or UPPER(search_type_code) like %s or start_time like %s  or start_time like %s  or stop_time like %s  or stop_latitude like %s  or stop_longitude like %s or stop_longitude like %s or se_fid like %s or ds_id like %s or site_fid like %s or temp like %s or conductivity like %s", filter, filter, filter, filter, filter, filter, filter, filter, filter, filter, filter, filter, filter, filter, filter)
	}

	countrows, err := s.db.Queryx(searchDataSummaryCountSql + filterQuery)
	if err != nil {
		return searchSummariesWithCount, err
	}

	for countrows.Next() {
		err = countrows.Scan(&searchSummariesWithCount.TotalCount)
		if err != nil {
			return searchSummariesWithCount, err
		}
	}
	defer countrows.Close()
	searchSummaries := []models.SearchSummary{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "se_id"
	}

	searchDataSummarySqlWithSearch := searchDataSummarySql + filterQuery + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))

	rows, err := s.db.Queryx(searchDataSummarySqlWithSearch)
	if err != nil {
		return searchSummariesWithCount, err
	}

	for rows.Next() {
		summary := models.SearchSummary{}
		err = rows.Scan(&summary.SeID, &summary.SearchDate, &summary.Recorder, &summary.SearchTypeCode, &summary.StartTime,
			&summary.StartLatitude, &summary.StartLongitude, &summary.StopTime, &summary.StopLatitude, &summary.StopLongitude, &summary.SeFID,
			&summary.DsID, &summary.SiteFID, &summary.Temp, &summary.Conductivity)
		if err != nil {
			return searchSummariesWithCount, err
		}
		searchSummaries = append(searchSummaries, summary)
	}

	defer rows.Close()
	searchSummariesWithCount.Items = searchSummaries

	return searchSummariesWithCount, err
}

var telemetryDataSummaryFullDataSql = `select * FROM table (pallid_data_api.telemetry_datasummary_fnc(:1, :2, :3, :4, :5, :6, :7, to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

func (s *PallidSturgeonStore) GetFullTelemetryDataSummary(year string, officeCode string, project string, approved string, season string, spice string, month string, fromDate string, toDate string) (string, error) {
	dbQuery, err := s.db.Prepare(telemetryDataSummaryFullDataSql)
	if err != nil {
		return "Cannot create file", err
	}

	rows, err := dbQuery.Query(year, officeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return "Cannot create file", err
	}
	defer rows.Close()

	cols, _ := rows.Columns()

	file, err := os.Create("TelemetryDataSummary.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	//save header
	data := make([]string, 0)
	data = append(data, cols...)
	err = writer.Write(data)
	if err != nil {
		log.Fatal("Cannot write to file", err)
	}

	for rows.Next() {

		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		data := make([]string, 0)
		for i := range cols {
			var v string
			val := columns[i]

			if val == nil {
				v = ""
			} else {
				v = fmt.Sprintf("%v", val)
			}
			data = append(data, v)
		}

		err := writer.Write(data)
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}

	return file.Name(), err
}

var telemetryDataSummarySql = `select t_id, 
	COALESCE(year, 0) as year, 
	COALESCE(field_office_code, 'ZZ') as field_office_code, 
	COALESCE(project_code, 0) as project_code, 
	COALESCE(segment_code, 0) as segment_code,
	COALESCE(season_code, '') as season_code,
	COALESCE(bend_number, 0) as bend_number, 
	radio_tag_num, 
	frequency_id, 
	capture_time, 
	capture_latitude, 
	capture_longitude, 
	COALESCE(position_confidence, 0) as position_confidence, 
	COALESCE(macro_code, '') as macro_code, 
	COALESCE(meso_code, '') as meso_code, 
	COALESCE(depth, 0) as depth, 
	COALESCE(conductivity, 0) as conductivity, 
	COALESCE(turbidity, 0) as turbidity FROM table (pallid_data_api.telemetry_datasummary_fnc(:1, :2, :3, :4, :5, :6, :7, to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

var telemetryDataSummaryCountSql = `select count(*) FROM table (pallid_data_api.telemetry_datasummary_fnc(:1, :2, :3, :4, :5, :6, :7, to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

func (s *PallidSturgeonStore) GetTelemetryDataSummary(year string, officeCode string, project string, approved string, season string, spice string, month string, fromDate string, toDate string, queryParams models.SearchParams) (models.TelemetrySummaryWithCount, error) {
	telemetrySummaryWithCount := models.TelemetrySummaryWithCount{}
	countQuery, err := s.db.Prepare(telemetryDataSummaryCountSql)
	if err != nil {
		return telemetrySummaryWithCount, err
	}

	countrows, err := countQuery.Query(year, officeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return telemetrySummaryWithCount, err
	}
	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&telemetrySummaryWithCount.TotalCount)
		if err != nil {
			return telemetrySummaryWithCount, err
		}
	}

	telemetrySummaries := []models.TelemetrySummary{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "t_id"
	}
	telemetryDataEntriesSqlWithSearch := telemetryDataSummarySql + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))

	dbQuery, err := s.db.Prepare(telemetryDataEntriesSqlWithSearch)
	if err != nil {
		return telemetrySummaryWithCount, err
	}

	rows, err := dbQuery.Query(year, officeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return telemetrySummaryWithCount, err
	}
	defer rows.Close()

	for rows.Next() {
		summary := models.TelemetrySummary{}
		err = rows.Scan(&summary.TId,
			&summary.Year,
			&summary.FieldOffice,
			&summary.Project,
			&summary.Segment,
			&summary.Season,
			&summary.Bend,
			&summary.RadioTagNum,
			&summary.FrequencyIdCode,
			&summary.CaptureTime,
			&summary.CaptureLatitude,
			&summary.CaptureLongitude,
			&summary.PositionConfidence,
			&summary.MacroId,
			&summary.MesoId,
			&summary.Depth,
			&summary.Conductivity,
			&summary.Turbidity)
		if err != nil {
			return telemetrySummaryWithCount, err
		}
		telemetrySummaries = append(telemetrySummaries, summary)
	}

	telemetrySummaryWithCount.Items = telemetrySummaries

	return telemetrySummaryWithCount, err
}

var procedureDataSummaryFullDataSql = `select * FROM table (pallid_data_api.procedure_datasummary_fnc(:1, :2, :3, :4, :5, :6, :7, to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

func (s *PallidSturgeonStore) GetFullProcedureDataSummary(year string, officeCode string, project string, approved string, season string, spice string, month string, fromDate string, toDate string) (string, error) {
	dbQuery, err := s.db.Prepare(procedureDataSummaryFullDataSql)
	if err != nil {
		return "Cannot create file", err
	}

	rows, err := dbQuery.Query(year, officeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return "Cannot create file", err
	}
	defer rows.Close()

	cols, _ := rows.Columns()

	file, err := os.Create("ProcedureDataSummary.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	//save header
	data := make([]string, 0)
	data = append(data, cols...)
	err = writer.Write(data)
	if err != nil {
		log.Fatal("Cannot write to file", err)
	}

	for rows.Next() {

		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		data := make([]string, 0)

		for i := range cols {
			var v string
			val := columns[i]

			if val == nil {
				v = ""
			} else {
				v = fmt.Sprintf("%v", val)
			}
			data = append(data, v)
		}

		err := writer.Write(data)
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}

	return file.Name(), err
}

var procedureDataSummarySql = `select pid_display, 
	mr_id, 
	COALESCE(year, 0) as year, 
	COALESCE(field_office_code, 'ZZ') as field_office_code, 
	COALESCE(project_code, 0) as project_code, 
	COALESCE(segment_code, 0) as segment_code,
	COALESCE(season_code, '') as season_code,
	purpose_code, 
	procedure_date, 
	COALESCE(new_radio_tag_num, 0) as new_radio_tag_num, 
	COALESCE(new_frequency_id, 0) as new_frequency_id, 
	COALESCE(spawn_code, '') as spawn_code, 
	COALESCE(expected_spawn_year, 0) as expected_spawn_year FROM table (pallid_data_api.procedure_datasummary_fnc(:1, :2, :3, :4, :5, :6, :7, to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

var procedureDataSummaryCountSql = `select count(*) FROM table (pallid_data_api.procedure_datasummary_fnc(:1, :2, :3, :4, :5, :6, :7, to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

func (s *PallidSturgeonStore) GetProcedureDataSummary(year string, officeCode string, project string, approved string, season string, spice string, month string, fromDate string, toDate string, queryParams models.SearchParams) (models.ProcedureSummaryWithCount, error) {
	procedureSummaryWithCount := models.ProcedureSummaryWithCount{}
	countQuery, err := s.db.Prepare(procedureDataSummaryCountSql)
	if err != nil {
		return procedureSummaryWithCount, err
	}

	countrows, err := countQuery.Query(year, officeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return procedureSummaryWithCount, err
	}
	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&procedureSummaryWithCount.TotalCount)
		if err != nil {
			return procedureSummaryWithCount, err
		}
	}

	procedureSummaries := []models.ProcedureSummary{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "mr_id"
	}
	procedureDataEntriesSqlWithSearch := procedureDataSummarySql + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(procedureDataEntriesSqlWithSearch)
	if err != nil {
		return procedureSummaryWithCount, err
	}

	rows, err := dbQuery.Query(year, officeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return procedureSummaryWithCount, err
	}
	defer rows.Close()

	for rows.Next() {
		summary := models.ProcedureSummary{}
		err = rows.Scan(&summary.ID,
			&summary.UniqueID,
			&summary.Year,
			&summary.FieldOffice,
			&summary.Project,
			&summary.Segment,
			&summary.Season,
			&summary.PurposeCode,
			&summary.ProcedureDate,
			&summary.NewRadioTagNum,
			&summary.NewFrequencyId,
			&summary.SpawnCode,
			&summary.ExpectedSpawnYear)
		if err != nil {
			return procedureSummaryWithCount, err
		}
		procedureSummaries = append(procedureSummaries, summary)
	}

	procedureSummaryWithCount.Items = procedureSummaries

	return procedureSummaryWithCount, err
}

var nextUploadSessionIdSql = `SELECT upload_session_seq.nextval from dual`

func (s *PallidSturgeonStore) GetUploadSessionId() (int, error) {
	rows, err := s.db.Query(nextUploadSessionIdSql)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var nextUploadSessionId int
	for rows.Next() {
		rows.Scan(&nextUploadSessionId)
	}

	return nextUploadSessionId, err
}

var insertUploadSiteSql = `insert into upload_sites (site_id, site_fid, site_year, fieldoffice_id, 
	field_office, project_id, project, 
	segment_id, segment, season_id, season, bend, bendrn, bend_river_mile, comments, 
	edit_initials, last_updated, upload_session_id, uploaded_by, upload_filename)
	values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20)`

func (s *PallidSturgeonStore) SaveSiteUpload(uploadSite models.UploadSite) error {
	_, err := s.db.Exec(insertUploadSiteSql,
		uploadSite.SiteID,
		uploadSite.SiteFID,
		uploadSite.SiteYear,
		uploadSite.FieldofficeID,
		uploadSite.FieldOffice,
		uploadSite.ProjectId,
		uploadSite.Project,
		uploadSite.SegmentId,
		uploadSite.Segment,
		uploadSite.SeasonId,
		uploadSite.Season,
		uploadSite.Bend,
		uploadSite.Bendrn,
		uploadSite.BendRiverMile,
		uploadSite.Comments,
		uploadSite.EditInitials,
		uploadSite.LastUpdated,
		uploadSite.UploadSessionId,
		uploadSite.UploadedBy,
		uploadSite.UploadFilename,
	)

	return err
}

var insertFishUploadSql = `insert into upload_fish (site_id, f_fid, mr_fid, panelhook, bait, species, length, weight,
	fishcount, fin_curl, otolith, rayspine, scale, ftprefix, ftnum, ftmr,
	comments, edit_initials, last_updated, upload_session_id, uploaded_by, upload_filename)
	 values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22)`

func (s *PallidSturgeonStore) SaveFishUpload(uploadFish models.UploadFish) error {
	_, err := s.db.Exec(insertFishUploadSql,
		uploadFish.SiteID,
		uploadFish.Ffid,
		uploadFish.MrFid,
		uploadFish.Panelhook,
		uploadFish.Bait,
		uploadFish.Species,
		uploadFish.Length,
		uploadFish.Weight,
		uploadFish.Fishcount,
		uploadFish.FinCurl,
		uploadFish.Otolith,
		uploadFish.Rayspine,
		uploadFish.Scale,
		uploadFish.Ftprefix,
		uploadFish.Ftnum,
		uploadFish.Ftmr,
		uploadFish.Comments,
		uploadFish.EditInitials,
		uploadFish.LastUpdated,
		uploadFish.UploadSessionId,
		uploadFish.UploadedBy,
		uploadFish.UploadFilename,
	)

	return err
}

var insertSearchUploadSql = `insert into upload_search(se_fid, ds_id, site_id, site_fid, search_date, recorder, search_type_code, search_day, start_time,  
		start_latitude, start_longitude, stop_time, stop_latitude, stop_longitude, temp, conductivity, edit_initials, last_updated, 
		upload_session_id, uploaded_by, upload_filename)
	values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21)`

func (s *PallidSturgeonStore) SaveSearchUpload(uploadSearch models.UploadSearch) error {
	_, err := s.db.Exec(insertSearchUploadSql,
		uploadSearch.SeFid,
		uploadSearch.DsId,
		uploadSearch.SiteID,
		uploadSearch.SiteFid,
		uploadSearch.SearchDateTime,
		uploadSearch.Recorder,
		uploadSearch.SearchTypeCode,
		uploadSearch.SearchDay,
		uploadSearch.StartTime,
		uploadSearch.StartLatitude,
		uploadSearch.StartLongitude,
		uploadSearch.StopTime,
		uploadSearch.StopLatitude,
		uploadSearch.StopLongitude,
		uploadSearch.Temp,
		uploadSearch.Conductivity,
		uploadSearch.EditInitials,
		uploadSearch.LastUpdated,
		uploadSearch.UploadSessionId,
		uploadSearch.UploadedBy,
		uploadSearch.UploadFilename,
	)

	return err
}

var insertSupplementalUploadSql = `insert into upload_supplemental (site_id, f_fid, mr_fid,
	tagnumber, pitrn,
	scuteloc, scutenum, scuteloc2, scutenum2,
	elhv, elcolor, erhv, ercolor, cwtyn, dangler, genetic, genetics_vial_number,
	broodstock, hatch_wild, species_id, archive,
	head, snouttomouth, inter, mouthwidth, m_ib,
	l_ob, l_ib, r_ib,
	r_ob, anal, dorsal, status, hatchery_origin,
	sex, stage, recapture, photo,
	genetic_needs, other_tag_info,
	comments, edit_initials,
	last_updated, upload_session_id, uploaded_by, upload_filename) values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,
		:21,:22,:23,:24,:25,:26,:27,:28,:29,:30,:31,:32,:33,:34,:35,:36,:37,:38,:39,:40,:41,:42,:43,:44,:45,:46)`

func (s *PallidSturgeonStore) SaveSupplementalUpload(uploadSupplemental models.UploadSupplemental) error {
	_, err := s.db.Exec(insertSupplementalUploadSql,
		uploadSupplemental.SiteID,
		uploadSupplemental.FFid,
		uploadSupplemental.MrFid,
		uploadSupplemental.Tagnumber,
		uploadSupplemental.Pitrn,
		uploadSupplemental.Scuteloc,
		uploadSupplemental.Scutenum,
		uploadSupplemental.Scuteloc2,
		uploadSupplemental.Scutenum2,
		uploadSupplemental.Elhv,
		uploadSupplemental.Elcolor,
		uploadSupplemental.Erhv,
		uploadSupplemental.Ercolor,
		uploadSupplemental.Cwtyn,
		uploadSupplemental.Dangler,
		uploadSupplemental.Genetic,
		uploadSupplemental.GeneticsVialNumber,
		uploadSupplemental.Broodstock,
		uploadSupplemental.HatchWild,
		uploadSupplemental.SpeciesId,
		uploadSupplemental.Archive,
		uploadSupplemental.Head,
		uploadSupplemental.Snouttomouth,
		uploadSupplemental.Inter,
		uploadSupplemental.Mouthwidth,
		uploadSupplemental.MIb,
		uploadSupplemental.LOb,
		uploadSupplemental.LIb,
		uploadSupplemental.RIb,
		uploadSupplemental.ROb,
		uploadSupplemental.Anal,
		uploadSupplemental.Dorsal,
		uploadSupplemental.Status,
		uploadSupplemental.HatcheryOrigin,
		uploadSupplemental.Sex,
		uploadSupplemental.Stage,
		uploadSupplemental.Recapture,
		uploadSupplemental.Photo,
		uploadSupplemental.GeneticNeeds,
		uploadSupplemental.OtherTagInfo,
		uploadSupplemental.Comments,
		uploadSupplemental.EditInitials,
		uploadSupplemental.LastUpdated,
		uploadSupplemental.UploadSessionId,
		uploadSupplemental.UploadedBy,
		uploadSupplemental.UploadFilename,
	)

	return err
}

var insertProcedureUploadSql = `insert into upload_procedure (f_fid, mr_fid, purpose_code, procedure_date, procedure_start_time, procedure_end_time, procedure_by, 
	antibiotic_injection_ind, photo_dorsal_ind, photo_ventral_ind, photo_left_ind,
	old_radio_tag_num, old_frequency_id, dst_serial_num, dst_start_date, dst_start_time, dst_reimplant_ind, new_radio_tag_num,
	new_frequency_id, sex_code, blood_sample_ind, egg_sample_ind, comments, fish_health_comments,
	eval_location_code, spawn_code, visual_repro_status_code, ultrasound_repro_status_code,
	expected_spawn_year, ultrasound_gonad_length, gonad_condition,
	edit_initials, last_updated, upload_session_id, uploaded_by, upload_filename)                                                        
values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23,:24,:25,:26,:27,:28,:29,:30,:31,:32,:33,:34,:35,:36)`

func (s *PallidSturgeonStore) SaveProcedureUpload(uploadProcedure models.UploadProcedure) error {
	_, err := s.db.Exec(insertProcedureUploadSql,
		uploadProcedure.FFid,
		uploadProcedure.MrFid,
		uploadProcedure.PurposeCode,
		uploadProcedure.ProcedureDateTime,
		uploadProcedure.ProcedureStartTime,
		uploadProcedure.ProcedureEndTime,
		uploadProcedure.ProcedureBy,
		uploadProcedure.AntibioticInjectionInd,
		uploadProcedure.PhotoDorsalInd,
		uploadProcedure.PhotoVentralInd,
		uploadProcedure.PhotoLeftInd,
		uploadProcedure.OldRadioTagNum,
		uploadProcedure.OldFrequencyId,
		uploadProcedure.DstSerialNum,
		uploadProcedure.DstStartDateTime,
		uploadProcedure.DstStartTime,
		uploadProcedure.DstReimplantInd,
		uploadProcedure.NewRadioTagNum,
		uploadProcedure.NewFrequencyId,
		uploadProcedure.SexCode,
		uploadProcedure.BloodSampleInd,
		uploadProcedure.EggSampleInd,
		uploadProcedure.Comments,
		uploadProcedure.FishHealthComments,
		uploadProcedure.EvalLocationCode,
		uploadProcedure.SpawnStatus,
		uploadProcedure.VisualReproStatusCode,
		uploadProcedure.UltrasoundReproStatusCode,
		uploadProcedure.ExpectedSpawnYear,
		uploadProcedure.UltrasoundGonadLength,
		uploadProcedure.GonadCondition,
		uploadProcedure.EditInitials,
		uploadProcedure.LastUpdated,
		uploadProcedure.UploadSessionId,
		uploadProcedure.UploadedBy,
		uploadProcedure.UploadFilename,
	)

	return err
}

var insertMoriverUploadSql = `insert into upload_mr (site_id, site_fid, mr_fid, season, setdate, subsample, subsamplepass, 
	subsamplen, recorder, gear, gear_type, temp, turbidity, conductivity, do,
	distance, width, netrivermile, structurenumber, usgs, riverstage, discharge,
	u1, u2, u3, u4, u5, u6, u7, macro, meso, habitatrn, qc,
	micro_structure, structure_flow, structure_mod, set_site_1, set_site_2, set_site_3,
	starttime, startlatitude, startlongitude, stoptime, stoplatitude, stoplongitude,
	depth1, velocitybot1, velocity08_1, velocity02or06_1,
	depth2, velocitybot2, velocity08_2, velocity02or06_2,
	depth3, velocitybot3, velocity08_3, velocity02or06_3,
	watervel, cobble, organic, silt, sand, gravel,
	comments, edit_initials, last_updated, upload_session_id,
	uploaded_by, upload_filename, complete,
	no_turbidity, no_velocity)

	 values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23,:24,
		:25,:26,:27,:28,:29,:30,:31,:32,:33,:34,:35,:36,:37,:38,:39,:40,:41,:42,:43,:44,:45,:46,:47,
		:48,:49,:50,:51,:52,:53,:54,:55,:56,:57,:58,:59,:60,:61,:62,:63,:64,:65,:66,:67,:68,:69,:70,
		:71,:72)`

func (s *PallidSturgeonStore) SaveMoriverUpload(UploadMoriver models.UploadMoriver) error {
	_, err := s.db.Exec(insertMoriverUploadSql,
		UploadMoriver.SiteID, UploadMoriver.SiteFid, UploadMoriver.MrFid, UploadMoriver.Season, UploadMoriver.SetdateTime,
		UploadMoriver.Subsample, UploadMoriver.Subsamplepass, UploadMoriver.Subsamplen, UploadMoriver.Recorder,
		UploadMoriver.Gear, UploadMoriver.GearType, UploadMoriver.Temp, UploadMoriver.Turbidity, UploadMoriver.Conductivity,
		UploadMoriver.Do, UploadMoriver.Distance, UploadMoriver.Width, UploadMoriver.Netrivermile, UploadMoriver.Structurenumber,
		UploadMoriver.Usgs, UploadMoriver.Riverstage, UploadMoriver.Discharge, UploadMoriver.U1, UploadMoriver.U2, UploadMoriver.U3, UploadMoriver.U4,
		UploadMoriver.U5, UploadMoriver.U6, UploadMoriver.U7, UploadMoriver.Macro, UploadMoriver.Meso, UploadMoriver.Habitatrn, UploadMoriver.Qc,
		UploadMoriver.MicroStructure, UploadMoriver.StructureFlow, UploadMoriver.StructureMod, UploadMoriver.SetSite1, UploadMoriver.SetSite2, UploadMoriver.SetSite3,
		UploadMoriver.StartTime, UploadMoriver.StartLatitude, UploadMoriver.StartLongitude, UploadMoriver.StopTime, UploadMoriver.StopLatitude, UploadMoriver.StopLongitude,
		UploadMoriver.Depth1, UploadMoriver.Velocitybot1, UploadMoriver.Velocity08_1, UploadMoriver.Velocity02or06_1,
		UploadMoriver.Depth2, UploadMoriver.Velocitybot2, UploadMoriver.Velocity08_2, UploadMoriver.Velocity02or06_2,
		UploadMoriver.Depth3, UploadMoriver.Velocitybot3, UploadMoriver.Velocity08_3, UploadMoriver.Velocity02or06_3,
		UploadMoriver.Watervel, UploadMoriver.Cobble, UploadMoriver.Organic, UploadMoriver.Silt, UploadMoriver.Sand, UploadMoriver.Gravel,
		UploadMoriver.Comments, UploadMoriver.EditInitials, UploadMoriver.LastUpdated, UploadMoriver.UploadSessionId,
		UploadMoriver.UploadedBy, UploadMoriver.UploadFilename, UploadMoriver.Complete,
		UploadMoriver.NoTurbidity, UploadMoriver.NoVelocity,
	)

	return err
}

var insertTelemetryUploadSql = `insert into upload_telemetry_fish(t_fid, se_fid, bend, radio_tag_num, frequency_id_code, capture_time, capture_latitude, capture_longitude,
	position_confidence, macro_id, meso_id, depth, temp, conductivity, turbidity, silt, sand, gravel, comments, edit_initials, last_updated, upload_session_id, uploaded_by, upload_filename)
	values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23,:24)`

func (s *PallidSturgeonStore) SaveTelemetryUpload(uploadTelemetry models.UploadTelemetry) error {
	_, err := s.db.Exec(insertTelemetryUploadSql,
		uploadTelemetry.TFid,
		uploadTelemetry.SeFid,
		uploadTelemetry.Bend,
		uploadTelemetry.RadioTagNum,
		uploadTelemetry.FrequencyIdCode,
		uploadTelemetry.CaptureTime,
		uploadTelemetry.CaptureLatitude,
		uploadTelemetry.CaptureLongitude,
		uploadTelemetry.PositionConfidence,
		uploadTelemetry.MacroId,
		uploadTelemetry.MesoId,
		uploadTelemetry.Depth,
		uploadTelemetry.Temp,
		uploadTelemetry.Conductivity,
		uploadTelemetry.Turbidity,
		uploadTelemetry.Silt,
		uploadTelemetry.Sand,
		uploadTelemetry.Gravel,
		uploadTelemetry.Comments,
		uploadTelemetry.EditInitials,
		uploadTelemetry.LastUpdated,
		uploadTelemetry.UploadSessionId,
		uploadTelemetry.UploadedBy,
		uploadTelemetry.UploadFilename,
	)

	return err
}

func (s *PallidSturgeonStore) CallStoreProcedures(uploadedBy string, uploadSessionId int) (models.ProcedureOut, error) {

	procedureOut := models.ProcedureOut{}

	uploadFishStmt, err := s.db.Prepare("begin PALLID_DATA_UPLOAD.uploadFinal (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12);  end;")
	if err != nil {
		return procedureOut, err
	}

	var p_site_cnt_final int
	var p_mr_cnt_final int
	var p_fishCntFinal int
	var p_searchCntFinal int
	var p_suppCntFinal int
	var p_telemetryCntFinal int
	var p_procedureCntFinal int
	var p_noSite_cnt int
	var p_siteMatch int
	var p_noSiteID_msg string

	_, err = uploadFishStmt.Exec(godror.PlSQLArrays, uploadedBy, sql.Out{Dest: &p_site_cnt_final}, sql.Out{Dest: &p_mr_cnt_final}, sql.Out{Dest: &p_fishCntFinal}, sql.Out{Dest: &p_searchCntFinal}, sql.Out{Dest: &p_suppCntFinal}, sql.Out{Dest: &p_telemetryCntFinal}, sql.Out{Dest: &p_procedureCntFinal}, sql.Out{Dest: &p_noSite_cnt}, sql.Out{Dest: &p_siteMatch}, sql.Out{Dest: &p_noSiteID_msg}, uploadSessionId)
	if err != nil {
		return procedureOut, err
	}

	uploadFishStmt.Close()

	procedureOut.UploadSessionId = uploadSessionId
	procedureOut.UploadedBy = uploadedBy
	procedureOut.SiteCntFinal = p_site_cnt_final
	procedureOut.MrCntFinal = p_mr_cnt_final
	procedureOut.FishCntFinal = p_fishCntFinal
	procedureOut.SearchCntFinal = p_searchCntFinal
	procedureOut.SuppCntFinal = p_suppCntFinal
	procedureOut.TelemetryCntFinal = p_telemetryCntFinal
	procedureOut.ProcedureCntFinal = p_procedureCntFinal
	procedureOut.NoSiteCnt = p_noSite_cnt
	procedureOut.SiteMatch = p_siteMatch
	procedureOut.NoSiteIDMsg = p_noSiteID_msg

	return procedureOut, err
}

var errorCountSql = `select el.year,count(el.el_id)
						from site_error_log_v el
						where NVL(error_fixed,0) = 0
						and  (case
						when el.worksheet_type_id = 2 then 
							(select FIELDOFFICE
							from ds_sites
							where site_id = (select site_id
											from ds_moriver
											where mr_id = el.worksheet_id))
						when el.worksheet_type_id = 1 then
							NULL
						when el.worksheet_type_id in(3,4) then
							(select FIELDOFFICE
							from ds_sites
							where site_id = (select mr2.site_id
											from ds_sites s2, ds_moriver mr2, ds_fish f2
											where s2.site_id = mr2.site_id
											and mr2.mr_id = F2.MR_ID
											and f2.f_id = el.worksheet_id))
						end) = :1
						group by el.year
						Order By el.year desc`

func (s *PallidSturgeonStore) GetErrorCount(fieldOfficeCode string) ([]models.ErrorCount, error) {
	errorCounts := []models.ErrorCount{}

	selectQuery, err := s.db.Prepare(errorCountSql)
	if err != nil {
		return errorCounts, err
	}

	rows, err := selectQuery.Query(fieldOfficeCode)
	if err != nil {
		return errorCounts, err
	}
	defer rows.Close()

	for rows.Next() {
		errorCount := models.ErrorCount{}
		err = rows.Scan(&errorCount.Year, &errorCount.Count)
		if err != nil {
			return errorCounts, err
		}
		errorCounts = append(errorCounts, errorCount)
	}

	return errorCounts, err
}

var usgNoVialNumberSql = `select fo.field_office_description||' : '||p.project_description as fp,
							f.species, 
							f.f_id, mr.mr_id, MR.SITE_ID as mrsite_id,   DS.SITE_ID as s_site_id,
							f.f_fid, Sup.GENETICS_VIAL_NUMBER

							from ds_fish f, ds_supplemental sup, ds_moriver mr, ds_sites ds,
							project_lk p, segment_lk s, field_office_lk fo

							where F.F_ID = Sup.F_ID (+)
							and MR.MR_ID = F.MR_ID (+)
							and mr.site_id = ds.site_id (+)
							and DS.PROJECT_ID = P.PROJECT_CODE (+)
							and DS.FIELDOFFICE = fo.FIELD_OFFICE_CODE (+)
							and ds.SEGMENT_ID = s.segment_code (+)

							and f.species = 'USG'
							and Sup.GENETICS_VIAL_NUMBER IS NULL
							and (CASE when :1 != 'ZZ' THEN
									ds.FIELDOFFICE
								ELSE
									:2
								END) = :3
							and ds.PROJECT_ID != 2
							order by ds.FIELDOFFICE, ds.PROJECT_ID, ds.SEGMENT_ID, ds.BEND`

func (s *PallidSturgeonStore) GetUsgNoVialNumbers(fieldOfficeCode string) ([]models.UsgNoVialNumber, error) {
	usgNoVialNumbers := []models.UsgNoVialNumber{}

	selectQuery, err := s.db.Prepare(usgNoVialNumberSql)
	if err != nil {
		return usgNoVialNumbers, err
	}

	rows, err := selectQuery.Query(fieldOfficeCode, fieldOfficeCode, fieldOfficeCode)
	if err != nil {
		return usgNoVialNumbers, err
	}
	defer rows.Close()

	for rows.Next() {
		usgNoVialNumber := models.UsgNoVialNumber{}
		err = rows.Scan(&usgNoVialNumber.Fp, &usgNoVialNumber.SpeciesCode, &usgNoVialNumber.FID, &usgNoVialNumber.MrID, &usgNoVialNumber.MrsiteID, &usgNoVialNumber.SSiteID, &usgNoVialNumber.FFID, &usgNoVialNumber.GeneticsVialNumber)
		if err != nil {
			return usgNoVialNumbers, err
		}
		usgNoVialNumbers = append(usgNoVialNumbers, usgNoVialNumber)
	}

	return usgNoVialNumbers, err
}

var unapprovedDataSheetsSql = `select 
								asv.ch,
								f.field_office_description||' : '||p.project_description as fp, 
								s.segment_description,
								m.BEND,
								m.MR_ID, 
								m.UNIQUEIDENTIFIER,
								m.SETDATE,
								m.SUBSAMPLE,
								m.RECORDER,
								m.CHECKBY,
								m.NETRIVERMILE,
								m.site_id,
								ds.PROJECT_ID, ds.SEGMENT_ID, ds.SEASON, ds.FIELDOFFICE,
								ds.SAMPLE_UNIT_TYPE,
								m.gear
								from DS_MORIVER m, project_lk p, segment_lk s, field_office_lk f, approval_status_v asv, ds_sites ds
								where m.site_id = ds.site_id (+)
								and ds.SEGMENT_ID = s.segment_code (+)
								and DS.PROJECT_ID = P.project_code (+)
								and DS.FIELDOFFICE = F.FIELD_OFFICE_CODE
								and m.mr_id = asv.mr_id (+)
								and asv.ch = 'Unapproved'
								-- and m.checkby is not null
								and asv.cb = 'YES'
								-- and asv.co = 'Complete'  
								-- and nvl(m.complete,0) = 1
								-- and nvl(m.approved,0) = 0
								and ds.PROJECT_ID != 2
								order by ds.FIELDOFFICE, ds.PROJECT_ID, ds.SEGMENT_ID, ds.BEND`

var unapprovedDataSheetsCountSql = `select 
							count(*)
							from DS_MORIVER m, project_lk p, segment_lk s, field_office_lk f, approval_status_v asv, ds_sites ds
							where m.site_id = ds.site_id (+)
							and ds.SEGMENT_ID = s.segment_code (+)
							and DS.PROJECT_ID = P.project_code (+)
							and DS.FIELDOFFICE = F.FIELD_OFFICE_CODE
							and m.mr_id = asv.mr_id (+)
							and asv.ch = 'Unapproved'
							-- and m.checkby is not null
							and asv.cb = 'YES'
							-- and asv.co = 'Complete'  
							-- and nvl(m.complete,0) = 1
							-- and nvl(m.approved,0) = 0
							and ds.PROJECT_ID != 2
							order by ds.FIELDOFFICE, ds.PROJECT_ID, ds.SEGMENT_ID, ds.BEND`

func (s *PallidSturgeonStore) GetUnapprovedDataSheets() (models.SummaryWithCount, error) {

	unapprovedDataSheetsWithCount := models.SummaryWithCount{}

	countrows, err := s.db.Query(unapprovedDataSheetsCountSql)
	if err != nil {
		return unapprovedDataSheetsWithCount, err
	}
	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&unapprovedDataSheetsWithCount.TotalCount)
		if err != nil {
			return unapprovedDataSheetsWithCount, err
		}
	}

	unapprovedDataSheets := make([]map[string]string, 0)

	rows, err := s.db.Query(unapprovedDataSheetsSql)
	if err != nil {
		return unapprovedDataSheetsWithCount, err
	}
	defer rows.Close()

	cols, _ := rows.Columns()

	for rows.Next() {

		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		data := make(map[string]string)

		for i, colName := range cols {
			var v string
			val := columns[i]

			if val == nil {
				v = ""
			} else {
				v = fmt.Sprintf("%v", val)
			}

			words := strings.Split(strings.ToLower(colName), "_")
			var convertedColName = words[0]
			if len(words) > 1 {
				word2 := strings.ToUpper(string(words[1][0])) + words[1][1:]
				convertedColName = words[0] + word2
			}

			data[convertedColName] = v
		}

		unapprovedDataSheets = append(unapprovedDataSheets, data)
	}

	unapprovedDataSheetsWithCount.Items = unapprovedDataSheets

	return unapprovedDataSheetsWithCount, err
}

var uncheckedDataSheetsSql = `select 
								asv.cb,
								p.project_description||' : '||s.segment_description||' : Bend '||ds.BEND as psb,
								m.MR_ID, 
								m.UNIQUEIDENTIFIER,
								m.SETDATE,
								m.SUBSAMPLE,
								m.RECORDER,
								m.CHECKBY,
								m.NETRIVERMILE,
								m.site_id,
								ds.PROJECT_ID, ds.SEGMENT_ID, ds.SEASON, ds.FIELDOFFICE, m.gear
								from DS_MORIVER m, project_lk p, segment_lk s, approval_status_v asv, ds_sites ds
								where m.site_id = ds.site_id (+)
								and ds.SEGMENT_ID = s.segment_code (+)
								and DS.PROJECT_ID = P.project_code (+)
								and m.mr_id = asv.mr_id (+)  
								and ds.FIELDOFFICE = :1 
								and asv.cb = 'Unchecked'
								-- and asv.co = 'Complete'
								and ds.PROJECT_ID != 2
								and M.MR_ID NOT IN (SELECT MR_ID 
													FROM DS_FISH
													WHERE SPECIES = 'BAFI')
								and m.FIELDOFFICE is not null`

var uncheckedDataSheetsCountSql = `select 
									count(*)
									from DS_MORIVER m, project_lk p, segment_lk s, approval_status_v asv, ds_sites ds
									where m.site_id = ds.site_id (+)
									and ds.segment_id = s.segment_code (+)
									and DS.PROJECT_ID = p.project_code (+)
									and m.mr_id = asv.mr_id (+)  
									and ds.FIELDOFFICE = :1 
									and asv.cb = 'Unchecked'
									-- and asv.co = 'Complete'
									and ds.PROJECT_ID != 2
									and M.MR_ID NOT IN (SELECT MR_ID 
														FROM DS_FISH
														WHERE SPECIES = 'BAFI')
									and m.FIELDOFFICE is not null`

func (s *PallidSturgeonStore) GetUncheckedDataSheets(fieldOfficeCode string, queryParams models.SearchParams) (models.SummaryWithCount, error) {

	uncheckedDataSheetsWithCount := models.SummaryWithCount{}
	countQuery, err := s.db.Prepare(uncheckedDataSheetsCountSql)
	if err != nil {
		return uncheckedDataSheetsWithCount, err
	}

	countrows, err := countQuery.Query(fieldOfficeCode)
	if err != nil {
		return uncheckedDataSheetsWithCount, err
	}
	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&uncheckedDataSheetsWithCount.TotalCount)
		if err != nil {
			return uncheckedDataSheetsWithCount, err
		}
	}

	uncheckedDataSheets := make([]map[string]string, 0)

	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "PROJECT_ID"
	}
	selectQueryWithSearch := uncheckedDataSheetsSql + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(selectQueryWithSearch)
	if err != nil {
		return uncheckedDataSheetsWithCount, err
	}

	rows, err := dbQuery.Query(fieldOfficeCode)
	if err != nil {
		return uncheckedDataSheetsWithCount, err
	}
	defer rows.Close()

	cols, _ := rows.Columns()

	for rows.Next() {

		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		data := make(map[string]string)

		for i, colName := range cols {
			var v string
			val := columns[i]

			if val == nil {
				v = ""
			} else {
				v = fmt.Sprintf("%v", val)
			}

			words := strings.Split(strings.ToLower(colName), "_")
			var convertedColName = words[0]
			if len(words) > 1 {
				word2 := strings.ToUpper(string(words[1][0])) + words[1][1:]
				convertedColName = words[0] + word2
			}

			data[convertedColName] = v
		}

		uncheckedDataSheets = append(uncheckedDataSheets, data)
	}

	uncheckedDataSheetsWithCount.Items = uncheckedDataSheets

	return uncheckedDataSheetsWithCount, err
}

func (s *PallidSturgeonStore) GetDownloadZip() (string, error) {

	rows, err := s.db.Query("SELECT content FROM media_tbl where md_id in (select max(md_id) from media_tbl)")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var data []byte
	for rows.Next() {
		rows.Scan(&data)
	}

	downloadInfo, err := s.GetDownloadInfo()
	if err != nil {
		return "", err
	}

	file, err := os.OpenFile(
		downloadInfo.Name,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal("Cannot write to file", err)
	}
	defer file.Close()

	bytesWritten, err := file.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes.\n", bytesWritten)
	return file.Name(), err

}

var uploadDownloadInfoSql = `insert into media_tbl (md_id, name, display_name, mime_type, content, last_updated) values ((select max(md_id)+1 from media_tbl),:1,:2,:3,:4,:5) returning md_id into :6`

func (s *PallidSturgeonStore) UploadDownloadZip(file *multipart.FileHeader) (int, error) {
	var id int
	fileContent, _ := file.Open()
	byteContainer, err := ioutil.ReadAll(fileContent)
	if err != nil {
		log.Fatal(err)
	}
	last5 := file.Filename[len(file.Filename)-9:]
	words := strings.Split(last5, "_")
	numbers := strings.Split(words[2], ".")
	version := "Version " + words[0] + "." + words[1] + "." + numbers[0]
	lastUpdated := time.Now()
	_, err = s.db.Exec(uploadDownloadInfoSql, file.Filename, version, "application/x-zip-compressed", byteContainer, lastUpdated, sql.Out{Dest: &id})

	return id, err
}

func (s *PallidSturgeonStore) GetDownloadInfo() (models.DownloadInfo, error) {
	downloadInfo := models.DownloadInfo{}

	rows, err := s.db.Query("SELECT name, display_name, last_updated FROM media_tbl where md_id in (select max(md_id) from media_tbl)")
	if err != nil {
		return downloadInfo, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&downloadInfo.Name, &downloadInfo.DisplayName, &downloadInfo.LastUpdated)
	}

	return downloadInfo, err
}
