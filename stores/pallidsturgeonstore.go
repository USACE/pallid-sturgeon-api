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

var getUserSql = `select u.id, u.username, u.first_name, u.last_name, u.email, r.description, f.FIELD_OFFICE_CODE, project_code from users_t u
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
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.OfficeCode, &user.ProjectCode)
		if err != nil {
			return user, err
		}
	}

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
	roles := []models.Role{}
	rows, err := s.db.Query("select * from role_lk order by id")
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

var getFieldOfficesSql = `select FO_ID, FIELD_OFFICE_CODE, FIELD_OFFICE_DESCRIPTION, STATE from field_office_lk where field_office_code <> 'ZZ' order by FO_ID`

var getFieldOfficeSqlAll = `select FO_ID, FIELD_OFFICE_CODE, FIELD_OFFICE_DESCRIPTION, STATE from field_office_lk order by FO_ID`

func (s *PallidSturgeonStore) GetFieldOffices(showAll string) ([]models.FieldOffice, error) {
	query := ""

	if showAll == "true" {
		query = getFieldOfficeSqlAll
	} else {
		query = getFieldOfficesSql
	}

	fieldOffices := []models.FieldOffice{}
	rows, err := s.db.Query(query)
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

func (s *PallidSturgeonStore) GetSampleUnitTypes() ([]models.SampleUnitType, error) {
	rows, err := s.db.Query("select * from sample_unit_type_lk where sample_unit_type_code <> 'A' order by SAMPLE_UNIT_TYPE_CODE")

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

func (s *PallidSturgeonStore) GetSeasons(projectCode string) ([]models.Season, error) {
	rows, err := s.db.Query("select s_id, season_code, season_description, field_app, project_code from season_lk where project_code = :1 order by s_id", projectCode)

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

var getSegmentsSql = `select segment_code as code, segment_description as description from fieldoffice_segment_v where field_office_code = :1 and project_code = :2`

func (s *PallidSturgeonStore) GetSegments(officeCode string, projectCode string) ([]models.Segment, error) {
	segments := []models.Segment{}

	selectQuery, err := s.db.Prepare(getSegmentsSql)
	if err != nil {
		return segments, err
	}

	rows, err := selectQuery.Query(officeCode, projectCode)
	if err != nil {
		return segments, err
	}
	defer rows.Close()

	for rows.Next() {
		segment := models.Segment{}
		err = rows.Scan(&segment.Code, &segment.Description)
		if err != nil {
			return nil, err
		}
		segments = append(segments, segment)
	}

	return segments, err
}

var getSampleUnitSql = `select sample_unit, sample_unit_desc as description from segment_sampleunit_v where segment_code = :1 and sample_unit_type = :2 order by 1`

func (s *PallidSturgeonStore) GetSampleUnit(sampleUnitType string, segmentCode string) ([]models.SampleUnit, error) {
	sampleUnits := []models.SampleUnit{}

	selectQuery, err := s.db.Prepare(getSampleUnitSql)
	if err != nil {
		return sampleUnits, err
	}

	rows, err := selectQuery.Query(segmentCode, sampleUnitType)
	if err != nil {
		return sampleUnits, err
	}
	defer rows.Close()

	for rows.Next() {
		sampleUnit := models.SampleUnit{}
		err = rows.Scan(&sampleUnit.SampleUnit, &sampleUnit.Description)
		if err != nil {
			return nil, err
		}
		sampleUnits = append(sampleUnits, sampleUnit)
	}

	return sampleUnits, err
}

func (s *PallidSturgeonStore) GetBendRn() ([]models.BendRn, error) {
	rows, err := s.db.Query("select bs_id, BEND_SELECTION_DESCRIPTION, BEND_SELECTION_CODE from BEND_SELECTION_LK order by 1 desc")

	bendRnsItems := []models.BendRn{}
	if err != nil {
		return bendRnsItems, err
	}
	defer rows.Close()

	for rows.Next() {
		bendRn := models.BendRn{}
		err = rows.Scan(&bendRn.ID, &bendRn.Description, &bendRn.Code)
		if err != nil {
			return nil, err
		}
		bendRnsItems = append(bendRnsItems, bendRn)
	}

	return bendRnsItems, err
}

func (s *PallidSturgeonStore) GetMeso(macro string) ([]models.Meso, error) {
	rows, err := s.db.Query("select MESOHABITAT_CODE from macro_meso_lk where MACROHABITAT_CODE = :1 order by 1 asc", macro)

	mesoItems := []models.Meso{}
	if err != nil {
		return mesoItems, err
	}
	defer rows.Close()

	for rows.Next() {
		meso := models.Meso{}
		err = rows.Scan(&meso.Code)
		if err != nil {
			return nil, err
		}
		mesoItems = append(mesoItems, meso)
	}

	return mesoItems, err
}

func (s *PallidSturgeonStore) GetStructureFlow(microStructure string) ([]models.StructureFlow, error) {
	rows, err := s.db.Query("select structure_flow, structure_flow_code from micro_habitat_desc_lk where micro_structure_code = :1 group by structure_flow, structure_flow_code", microStructure)

	structureFlowItems := []models.StructureFlow{}
	if err != nil {
		return structureFlowItems, err
	}
	defer rows.Close()

	for rows.Next() {
		structureFlow := models.StructureFlow{}
		err = rows.Scan(&structureFlow.Code, &structureFlow.ID)
		if err != nil {
			return nil, err
		}
		structureFlowItems = append(structureFlowItems, structureFlow)
	}

	return structureFlowItems, err
}

func (s *PallidSturgeonStore) GetStructureMod(structureFlow string) ([]models.StructureMod, error) {
	rows, err := s.db.Query("select structure_mod, structure_mod_code from micro_habitat_desc_lk where structure_flow_code = :1 group by structure_mod, structure_mod_code", structureFlow)

	structureModItems := []models.StructureMod{}
	if err != nil {
		return structureModItems, err
	}
	defer rows.Close()

	for rows.Next() {
		structureMod := models.StructureMod{}
		err = rows.Scan(&structureMod.Description, &structureMod.Code)
		if err != nil {
			return nil, err
		}
		structureModItems = append(structureModItems, structureMod)
	}

	return structureModItems, err
}

func (s *PallidSturgeonStore) GetSpecies() ([]models.Species, error) {
	rows, err := s.db.Query("select alpha_code from fish_code_lk order by 1 asc")

	speciesItems := []models.Species{}
	if err != nil {
		return speciesItems, err
	}
	defer rows.Close()

	for rows.Next() {
		species := models.Species{}
		err = rows.Scan(&species.Code)
		if err != nil {
			return nil, err
		}
		speciesItems = append(speciesItems, species)
	}

	return speciesItems, err
}

func (s *PallidSturgeonStore) GetFtPrefixes() ([]models.FtPrefix, error) {
	rows, err := s.db.Query("select tag_prefix_code from floy_tag_prefix_code_lk order by 1 asc")

	ftPrefixItems := []models.FtPrefix{}
	if err != nil {
		return ftPrefixItems, err
	}
	defer rows.Close()

	for rows.Next() {
		ftPrefix := models.FtPrefix{}
		err = rows.Scan(&ftPrefix.Code)
		if err != nil {
			return nil, err
		}
		ftPrefixItems = append(ftPrefixItems, ftPrefix)
	}

	return ftPrefixItems, err
}

func (s *PallidSturgeonStore) GetMr() ([]models.Mr, error) {
	rows, err := s.db.Query("select mark_recapture_code, mark_recapture_description from mark_recapture_lk order by 1 asc")

	mrItems := []models.Mr{}
	if err != nil {
		return mrItems, err
	}
	defer rows.Close()

	for rows.Next() {
		mr := models.Mr{}
		err = rows.Scan(&mr.Code, &mr.Description)
		if err != nil {
			return nil, err
		}
		mrItems = append(mrItems, mr)
	}

	return mrItems, err
}

func (s *PallidSturgeonStore) GetOtolith() ([]models.Otolith, error) {
	rows, err := s.db.Query("select code, description from fish_structure_lk order by 1 asc")

	otolithItems := []models.Otolith{}
	if err != nil {
		return otolithItems, err
	}
	defer rows.Close()

	for rows.Next() {
		otolith := models.Otolith{}
		err = rows.Scan(&otolith.Code, &otolith.Description)
		if err != nil {
			return nil, err
		}
		otolithItems = append(otolithItems, otolith)
	}

	return otolithItems, err
}

func (s *PallidSturgeonStore) GetSetSite1(microstructure string) ([]models.SetSite1, error) {
	rows, err := s.db.Query("select set_site_1, set_site_1_code from micro_set_site_lk where structure_code=:1 group by set_site_1, set_site_1_code order by set_site_1_code asc", microstructure)

	setSiteItems := []models.SetSite1{}
	if err != nil {
		return setSiteItems, err
	}
	defer rows.Close()

	for rows.Next() {
		setSite := models.SetSite1{}
		err = rows.Scan(&setSite.Description, &setSite.Code)
		if err != nil {
			return nil, err
		}
		setSiteItems = append(setSiteItems, setSite)
	}

	return setSiteItems, err
}

func (s *PallidSturgeonStore) GetSetSite2(setsite1 string) ([]models.SetSite2, error) {
	rows, err := s.db.Query("select set_site_two, set_site_two_code from micro_set_site_lk where set_site_1_code=:1 and set_site_two != 'SIDE  NOTCH' group by set_site_two, set_site_two_code order by set_site_two_code asc", setsite1)

	setSiteItems := []models.SetSite2{}
	if err != nil {
		return setSiteItems, err
	}
	defer rows.Close()

	for rows.Next() {
		setSite := models.SetSite2{}
		err = rows.Scan(&setSite.Description, &setSite.Code)
		if err != nil {
			return nil, err
		}
		setSiteItems = append(setSiteItems, setSite)
	}

	return setSiteItems, err
}

var headerDataSql = `select si.site_id, si.year, si.fieldoffice, si.project_id, si.segment_id, si.season, si.bend,
si.bendrn, si.sample_unit_type, COALESCE(mo.bendrivermile,0.0) as bendrivermile from ds_sites si
inner join ds_moriver mo on si.site_id = mo.site_id
where si.site_id=:1
group by si.site_id, si.year, si.fieldoffice, si.project_id, si.segment_id, si.season, si.bend,
si.bendrn, si.sample_unit_type, mo.bendrivermile`

func (s *PallidSturgeonStore) GetHeaderData(siteId string) ([]models.HeaderData, error) {
	rows, err := s.db.Query(headerDataSql, siteId)

	headerDataItems := []models.HeaderData{}
	if err != nil {
		return headerDataItems, err
	}
	defer rows.Close()

	for rows.Next() {
		headerData := models.HeaderData{}
		err = rows.Scan(&headerData.SiteId, &headerData.Year, &headerData.FieldOffice, &headerData.Project, &headerData.Segment, &headerData.Season,
			&headerData.Bend, &headerData.Bendrn, &headerData.SampleUnitType, &headerData.BendRiverMile)
		if err != nil {
			return nil, err
		}
		headerDataItems = append(headerDataItems, headerData)
	}

	return headerDataItems, err
}

var siteDataEntriesSql = `select si.SITE_ID,si.YEAR,si.FIELDOFFICE,si.PROJECT_ID,si.SEGMENT_ID,si.SEASON,si.BEND,si.BENDRN,si.SITE_FID,si.UPLOADED_BY,si.LAST_EDIT_COMMENT,si.EDIT_INITIALS,si.complete,
si.approved,si.UPLOAD_FILENAME,si.UPLOAD_SESSION_ID,si.SAMPLE_UNIT_TYPE,si.brm_id,fnc.bkg_color,fnc.bend_river_mile
from ds_sites si inner join table (pallid_data_entry_api.data_entry_site_fnc(:2,:3,:4,:5,:6,:7)) fnc on si.site_id = fnc.site_id`

var siteDataEntriesCountSql = `SELECT count(*) from ds_sites si inner join table (pallid_data_entry_api.data_entry_site_fnc(:2,:3,:4,:5,:6,:7)) fnc on si.site_id = fnc.site_id`

var siteDataEntriesBySiteIdSql = `select si.SITE_ID,si.YEAR,si.FIELDOFFICE,si.PROJECT_ID,si.SEGMENT_ID,si.SEASON,si.BEND,si.BENDRN,si.SITE_FID,si.UPLOADED_BY,si.LAST_EDIT_COMMENT,si.EDIT_INITIALS,si.complete,
si.approved,si.UPLOAD_FILENAME,si.UPLOAD_SESSION_ID,si.SAMPLE_UNIT_TYPE,si.brm_id,fnc.bkg_color,fnc.bend_river_mile
from ds_sites si inner join table (pallid_data_entry_api.data_entry_site_fnc(:2,:3,:4,:5,:6,:7)) fnc on si.site_id = fnc.site_id where site_id=:1`

var siteDataEntriesCountBySiteIdSql = `SELECT count(*) from ds_sites si inner join table (pallid_data_entry_api.data_entry_site_fnc(:2,:3,:4,:5,:6,:7)) fnc on si.site_id = fnc.site_id where site_id=:1`

func (s *PallidSturgeonStore) GetSiteDataEntries(siteId string, year string, officeCode string, project string, segment string, season string, bend string, queryParams models.SearchParams) (models.SitesWithCount, error) {
	siteDataEntryWithCount := models.SitesWithCount{}
	query := ""
	queryWithCount := ""
	id := ""

	if siteId != "" {
		query = siteDataEntriesBySiteIdSql
		queryWithCount = siteDataEntriesCountBySiteIdSql
		id = siteId
	}

	if siteId == "" {
		query = siteDataEntriesSql
		queryWithCount = siteDataEntriesCountSql
	}

	countQuery, err := s.db.Prepare(queryWithCount)
	if err != nil {
		return siteDataEntryWithCount, err
	}

	var countrows *sql.Rows
	if id == "" {
		countrows, err = countQuery.Query(year, officeCode, project, bend, season, segment)
		if err != nil {
			return siteDataEntryWithCount, err
		}
	} else {
		countrows, err = countQuery.Query(year, officeCode, project, bend, season, segment, id)
		if err != nil {
			return siteDataEntryWithCount, err
		}
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
		queryParams.OrderBy = "site_id desc"
	}
	siteDataEntriesSqlWithSearch := query + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(siteDataEntriesSqlWithSearch)
	if err != nil {
		return siteDataEntryWithCount, err
	}

	var rows *sql.Rows
	if id == "" {
		rows, err = dbQuery.Query(year, officeCode, project, bend, season, segment)
		if err != nil {
			return siteDataEntryWithCount, err
		}
	} else {
		rows, err = dbQuery.Query(year, officeCode, project, bend, season, segment, id)
		if err != nil {
			return siteDataEntryWithCount, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		siteDataEntry := models.Sites{}
		err = rows.Scan(&siteDataEntry.SiteID, &siteDataEntry.Year, &siteDataEntry.FieldofficeId, &siteDataEntry.ProjectId, &siteDataEntry.SegmentId, &siteDataEntry.SeasonId, &siteDataEntry.Bend, &siteDataEntry.Bendrn, &siteDataEntry.SiteFID,
			&siteDataEntry.UploadedBy, &siteDataEntry.LastEditComment, &siteDataEntry.EditInitials, &siteDataEntry.Complete, &siteDataEntry.Approved, &siteDataEntry.UploadFilename, &siteDataEntry.UploadSessionId,
			&siteDataEntry.SampleUnitTypeCode, &siteDataEntry.BrmID, &siteDataEntry.BkgColor, &siteDataEntry.BendRiverMile)
		if err != nil {
			return siteDataEntryWithCount, err
		}
		siteEntries = append(siteEntries, siteDataEntry)
	}

	siteDataEntryWithCount.Items = siteEntries

	return siteDataEntryWithCount, err
}

var insertSiteDataSql = `insert into ds_sites (brm_id, site_fid, year, FIELDOFFICE, PROJECT_ID,
	SEGMENT_ID, SEASON, SAMPLE_UNIT_TYPE, bend, BENDRN, edit_initials, last_updated, last_edit_comment, uploaded_by) 
	values ((CASE 
	when :14 = 'B' THEN (select brm_id from bend_river_mile_lk where bend_num = :15 and b_segment = :16)
	when :17 = 'C' THEN (select chute_id from chute_lk where chute_code = :18 and segment_id = :19)
	when :20 = 'R' THEN (select reach_id from reach_lk where reach_code = :21 and segment_id = :22)
	ELSE 0
	END),:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13) returning site_id into :23`

func (s *PallidSturgeonStore) SaveSiteDataEntry(code string, sampleUnitType string, segmentCode string, sitehDataEntry models.Sites) (int, error) {
	var id int
	_, err := s.db.Exec(insertSiteDataSql, sampleUnitType, code, segmentCode, sampleUnitType, code, segmentCode, sampleUnitType, code, segmentCode, sitehDataEntry.SiteFID, sitehDataEntry.Year, sitehDataEntry.FieldofficeId, sitehDataEntry.ProjectId,
		sitehDataEntry.SegmentId, sitehDataEntry.SeasonId, sitehDataEntry.SampleUnitTypeCode, sitehDataEntry.Bend, sitehDataEntry.Bendrn, sitehDataEntry.EditInitials, sitehDataEntry.LastUpdated,
		sitehDataEntry.LastEditComment, sitehDataEntry.UploadedBy, sql.Out{Dest: &id})

	return id, err
}

var updateSiteDataSql = `UPDATE ds_sites
SET site_fid = :2,
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
	  brm_id = :13,
		last_edit_comment = :14
WHERE site_id = :1`

func (s *PallidSturgeonStore) UpdateSiteDataEntry(sitehDataEntry models.Sites) error {
	_, err := s.db.Exec(updateSiteDataSql, sitehDataEntry.SiteFID, sitehDataEntry.Year, sitehDataEntry.FieldofficeId, sitehDataEntry.ProjectId, sitehDataEntry.SegmentId, sitehDataEntry.SeasonId, sitehDataEntry.SampleUnitTypeCode,
		sitehDataEntry.Bendrn, sitehDataEntry.EditInitials, sitehDataEntry.LastUpdated, sitehDataEntry.UploadedBy, sitehDataEntry.BrmID, sitehDataEntry.LastEditComment, sitehDataEntry.SiteID)
	return err
}

var fishDataEntriesSql = `select fi.f_id, fi.f_fid, fi.mr_id, si.site_id, fi.panelhook,fi.bait,fi.species, fi.length, fi.weight, fi.fishcount, fi.otolith, fi.rayspine, fi.scale, fi.ftprefix, fi.ftnum, fi.ftmr, fi.edit_initials, 
fi.last_edit_comment, fi.uploaded_by, fi.genetics_vial_number, fi.condition, fi.fin_curl from ds_fish fi inner join ds_moriver mo on fi.mr_id = mo.mr_id inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :1 != 'ZZ' THEN si.fieldoffice ELSE :2 END) = :3`

var fishDataEntriesCountSql = `select count(*) from ds_fish fi inner join ds_moriver mo on fi.mr_id = mo.mr_id inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :1 != 'ZZ' THEN si.fieldoffice ELSE :2 END) = :3`

var fishDataEntriesByFidSql = `select fi.f_id, fi.f_fid, fi.mr_id, si.site_id, fi.panelhook,fi.bait,fi.species, fi.length, fi.weight, fi.fishcount, fi.otolith, fi.rayspine, fi.scale, fi.ftprefix, fi.ftnum, fi.ftmr, fi.edit_initials, 
fi.last_edit_comment, fi.uploaded_by, fi.genetics_vial_number, fi.condition, fi.fin_curl from ds_fish fi inner join ds_moriver mo on fi.mr_id = mo.mr_id inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4 
and fi.f_id = :1`

var fishDataEntriesCountByFidSql = `select count(*) from ds_fish fi inner join ds_moriver mo on fi.mr_id = mo.mr_id inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4 
and fi.f_id = :1`

var fishDataEntriesByFfidSql = `select fi.f_id, fi.f_fid, fi.mr_id, si.site_id, fi.panelhook,fi.bait,fi.species, fi.length, fi.weight, fi.fishcount, fi.otolith, fi.rayspine, fi.scale, fi.ftprefix, fi.ftnum, fi.ftmr, fi.edit_initials, 
fi.last_edit_comment, fi.uploaded_by, fi.genetics_vial_number, fi.condition, fi.fin_curl from ds_fish fi inner join ds_moriver mo on fi.mr_id = mo.mr_id inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4
and fi.f_fid = :1`

var fishDataEntriesCountByFfidSql = `select count(*) from ds_fish fi inner join ds_moriver mo on fi.mr_id = mo.mr_id inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4
and fi.f_fid = :1`

var fishDataEntriesByMridSql = `select fi.f_id, fi.f_fid, fi.mr_id, si.site_id, fi.panelhook,fi.bait,fi.species, fi.length, fi.weight, fi.fishcount, fi.otolith, fi.rayspine, fi.scale, fi.ftprefix, fi.ftnum, fi.ftmr, fi.edit_initials, 
fi.last_edit_comment, fi.uploaded_by, fi.genetics_vial_number, fi.condition, fi.fin_curl from ds_fish fi inner join ds_moriver mo on fi.mr_id = mo.mr_id inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4 
and fi.mr_id = :1`

var fishDataEntriesCountByMridSql = `select count(*) from ds_fish fi inner join ds_moriver mo on fi.mr_id = mo.mr_id inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4 
and fi.mr_id = :1`

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
		countrows, err = countQuery.Query(officeCode, officeCode, officeCode)
		if err != nil {
			return fishDataEntryWithCount, err
		}
	} else {
		countrows, err = countQuery.Query(officeCode, officeCode, officeCode, id)
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
		queryParams.OrderBy = "f_id asc"
	}
	fishDataEntriesSqlWithSearch := query + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(fishDataEntriesSqlWithSearch)
	if err != nil {
		return fishDataEntryWithCount, err
	}

	var rows *sql.Rows
	if id == "" {
		rows, err = dbQuery.Query(officeCode, officeCode, officeCode)
		if err != nil {
			return fishDataEntryWithCount, err
		}
	} else {
		rows, err = dbQuery.Query(officeCode, officeCode, officeCode, id)
		if err != nil {
			return fishDataEntryWithCount, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		fishDataEntry := models.UploadFish{}
		err = rows.Scan(&fishDataEntry.Fid, &fishDataEntry.Ffid, &fishDataEntry.MrID, &fishDataEntry.SiteID, &fishDataEntry.Panelhook, &fishDataEntry.Bait, &fishDataEntry.Species, &fishDataEntry.Length, &fishDataEntry.Weight, &fishDataEntry.Fishcount, &fishDataEntry.Otolith, &fishDataEntry.Rayspine,
			&fishDataEntry.Scale, &fishDataEntry.Ftprefix, &fishDataEntry.Ftnum, &fishDataEntry.Ftmr, &fishDataEntry.EditInitials, &fishDataEntry.LastEditComment, &fishDataEntry.UploadedBy, &fishDataEntry.GeneticsVialNumber, &fishDataEntry.Condition, &fishDataEntry.FinCurl)
		if err != nil {
			return fishDataEntryWithCount, err
		}
		fishEntries = append(fishEntries, fishDataEntry)
	}

	fishDataEntryWithCount.Items = fishEntries

	return fishDataEntryWithCount, err
}

var insertFishDataSql = `insert into ds_fish (FIELDOFFICE,PROJECT,SEGMENT,uniqueidentifier,id,panelhook,bait,SPECIES,length,weight,FISHCOUNT,otolith,rayspine,scale,FTPREFIX,FTNUM,FTMR,mr_id,edit_initials,last_edit_comment, 
last_updated, uploaded_by, genetics_vial_number, condition, fin_curl) 
values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23,:24,:25) returning f_id into :26`

func (s *PallidSturgeonStore) SaveFishDataEntry(fishDataEntry models.UploadFish) (int, error) {
	var id int
	_, err := s.db.Exec(insertFishDataSql, fishDataEntry.Fieldoffice, fishDataEntry.Project, fishDataEntry.Segment, fishDataEntry.UniqueID, fishDataEntry.Id, fishDataEntry.Panelhook,
		fishDataEntry.Bait, fishDataEntry.Species, fishDataEntry.Length, fishDataEntry.Weight, fishDataEntry.Fishcount, fishDataEntry.Otolith, fishDataEntry.Rayspine,
		fishDataEntry.Scale, fishDataEntry.Ftprefix, fishDataEntry.Ftnum, fishDataEntry.Ftmr, fishDataEntry.MrID, fishDataEntry.EditInitials, fishDataEntry.LastEditComment, fishDataEntry.LastUpdated, fishDataEntry.UploadedBy,
		fishDataEntry.GeneticsVialNumber, fishDataEntry.Condition, fishDataEntry.FinCurl, sql.Out{Dest: &id})

	return id, err
}

var updateFishDataSql = `UPDATE ds_fish SET
FIELDOFFICE = :2,
PROJECT = :3,
SEGMENT = :4,
uniqueidentifier = :5,
id = :6,
panelhook = :7,
bait = :8,
SPECIES = :9,
length = :10,
weight = :11,
FISHCOUNT = :12,
otolith = :13,
rayspine = :14,
scale = :15,
FTPREFIX = :16,
FTNUM = :17,
FTMR = :18,
edit_initials = :19,
last_edit_comment = :20,
last_updated = :21, 
uploaded_by = :22,
genetics_vial_number = :23,
condition = :24,
fin_curl = :25
WHERE f_id = :1`

func (s *PallidSturgeonStore) UpdateFishDataEntry(fishDataEntry models.UploadFish) error {
	_, err := s.db.Exec(updateFishDataSql, fishDataEntry.Fieldoffice, fishDataEntry.Project, fishDataEntry.Segment, fishDataEntry.UniqueID, fishDataEntry.Id, fishDataEntry.Panelhook,
		fishDataEntry.Bait, fishDataEntry.Species, fishDataEntry.Length, fishDataEntry.Weight, fishDataEntry.Fishcount, fishDataEntry.Otolith, fishDataEntry.Rayspine,
		fishDataEntry.Scale, fishDataEntry.Ftprefix, fishDataEntry.Ftnum, fishDataEntry.Ftmr, fishDataEntry.EditInitials, fishDataEntry.LastEditComment, fishDataEntry.LastUpdated, fishDataEntry.UploadedBy,
		fishDataEntry.GeneticsVialNumber, fishDataEntry.Condition, fishDataEntry.FinCurl, fishDataEntry.Fid)
	return err
}

func (s *PallidSturgeonStore) DeleteFishDataEntry(id string) error {
	_, err := s.db.Exec("delete from ds_fish where f_id = :1", id)
	return err
}

var insertMoriverDataSql = `insert into ds_moriver(mr_fid,site_id,FIELDOFFICE,PROJECT,SEGMENT,SEASON,setdate, subsample, subsamplepass, subsamplen, recorder, 
	gear, GEAR_TYPE, temp, turbidity, conductivity, do, distance, width, netrivermile, structurenumber, usgs, riverstage, discharge,
	u1, u2, u3, u4, u5, u6, u7, MACRO, MESO, habitatrn, qc, micro_structure, structure_flow, structure_mod, set_site_1, set_site_2, 
	set_site_3, starttime, startlatitude, startlongitude, stoptime, stoplatitude, stoplongitude, depth1, velocitybot1, velocity08_1, 
	velocity02or06_1, depth2, velocitybot2, velocity08_2, velocity02or06_2, depth3, velocitybot3, velocity08_3, velocity02or06_3, 
	watervel, cobble, ORGANIC, silt, sand, gravel, comments, complete, checkby, no_turbidity, no_velocity, edit_initials,last_edit_comment, 
	last_updated, uploaded_by, bend, bendrn, bendrivermile) values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,
		:21,:22,:23,:24,:25,:26,:27,:28,:29,:30,:31,:32,:33,:34,:35,:36,:37,:38,:39,:40,:41,:42,:43,:44,:45,:46,:47,:48,:49,:50,
		:51,:52,:53,:54,:55,:56,:57,:58,:59,:60,:61,:62,:63,:64,:65,:66,:67,:68,:69,:70,:71,:72,:73,:74,:75,:76,:77) returning mr_id into :78`

func (s *PallidSturgeonStore) SaveMoriverDataEntry(moriverDataEntry models.UploadMoriver) (int, error) {
	var id int
	_, err := s.db.Exec(insertMoriverDataSql, moriverDataEntry.MrFid, moriverDataEntry.SiteID, moriverDataEntry.FieldOffice,
		moriverDataEntry.Project, moriverDataEntry.Segment, moriverDataEntry.Season, moriverDataEntry.SetDate, moriverDataEntry.Subsample, moriverDataEntry.Subsamplepass,
		moriverDataEntry.Subsamplen, moriverDataEntry.Recorder, moriverDataEntry.Gear, moriverDataEntry.GearType, moriverDataEntry.Temp, moriverDataEntry.Turbidity, moriverDataEntry.Conductivity, moriverDataEntry.Do,
		moriverDataEntry.Distance, moriverDataEntry.Width, moriverDataEntry.Netrivermile, moriverDataEntry.Structurenumber, moriverDataEntry.Usgs, moriverDataEntry.Riverstage, moriverDataEntry.Discharge,
		moriverDataEntry.U1, moriverDataEntry.U2, moriverDataEntry.U3, moriverDataEntry.U4, moriverDataEntry.U5, moriverDataEntry.U6, moriverDataEntry.U7, moriverDataEntry.Macro, moriverDataEntry.Meso, moriverDataEntry.Habitatrn, moriverDataEntry.Qc,
		moriverDataEntry.MicroStructure, moriverDataEntry.StructureFlow, moriverDataEntry.StructureMod, moriverDataEntry.SetSite1, moriverDataEntry.SetSite2, moriverDataEntry.SetSite3,
		moriverDataEntry.StartTime, moriverDataEntry.StartLatitude, moriverDataEntry.StartLongitude, moriverDataEntry.StopTime, moriverDataEntry.StopLatitude, moriverDataEntry.StopLongitude,
		moriverDataEntry.Depth1, moriverDataEntry.Velocitybot1, moriverDataEntry.Velocity08_1, moriverDataEntry.Velocity02or06_1,
		moriverDataEntry.Depth2, moriverDataEntry.Velocitybot2, moriverDataEntry.Velocity08_2, moriverDataEntry.Velocity02or06_2,
		moriverDataEntry.Depth3, moriverDataEntry.Velocitybot3, moriverDataEntry.Velocity08_3, moriverDataEntry.Velocity02or06_3,
		moriverDataEntry.Watervel, moriverDataEntry.Cobble, moriverDataEntry.Organic, moriverDataEntry.Silt, moriverDataEntry.Sand, moriverDataEntry.Gravel,
		moriverDataEntry.Comments, moriverDataEntry.Complete, moriverDataEntry.Checkby, moriverDataEntry.NoTurbidity, moriverDataEntry.NoVelocity, moriverDataEntry.EditInitials, moriverDataEntry.LastEditComment, moriverDataEntry.LastUpdated, moriverDataEntry.UploadedBy,
		moriverDataEntry.Bend, moriverDataEntry.BendRn, moriverDataEntry.BendRiverMile, sql.Out{Dest: &id})
	return id, err
}

var updateMoriverDataSql = `UPDATE ds_moriver
SET  project = :2,segment = :3,SEASON = :4,setdate = :5, subsample = :6, subsamplepass = :7, subsamplen = :8, recorder = :9, gear = :10, 
GEAR_TYPE = :11, temp = :12, turbidity = :13, conductivity = :14, do = :15, distance = :16, width = :17, netrivermile = :18, structurenumber = :19, 
usgs = :20, riverstage = :21, discharge = :22, u1 = :23, u2 = :24, u3 = :25, u4 = :26, u5 = :27, u6 = :28, u7 = :29, MACRO = :30, MESO = :31, 
habitatrn = :32, qc = :33, micro_structure = :34, structure_flow = :35, structure_mod = :36, set_site_1 = :37, set_site_2 = :38, set_site_3 = :39,
starttime = :40, startlatitude = :41, startlongitude = :42, stoptime = :43, stoplatitude = :44, stoplongitude = :45, 
depth1 = :46, velocitybot1 = :47, velocity08_1 = :48, velocity02or06_1 = :49,
depth2 = :50, velocitybot2 = :51, velocity08_2 = :52, velocity02or06_2 = :53,
depth3 = :54, velocitybot3 = :55, velocity08_3 = :56, velocity02or06_3 = :57, 
watervel = :58, cobble = :59, ORGANIC = :60, silt = :61, sand = :62, gravel = :63, comments = :64, complete = :65, checkby = :66, 
no_turbidity = :67, no_velocity = :68, edit_initials = :69,  mr_fid= :70, site_id = :71, FIELDOFFICE = :72, last_edit_comment = :73, last_updated = :74, 
uploaded_by = :75, bend = :76, bendrn = :77, bendrivermile =:78 WHERE mr_id = :1`

func (s *PallidSturgeonStore) UpdateMoriverDataEntry(moriverDataEntry models.UploadMoriver) error {
	_, err := s.db.Exec(updateMoriverDataSql,
		moriverDataEntry.Project, moriverDataEntry.Segment, moriverDataEntry.Season, moriverDataEntry.SetDate, moriverDataEntry.Subsample, moriverDataEntry.Subsamplepass,
		moriverDataEntry.Subsamplen, moriverDataEntry.Recorder, moriverDataEntry.Gear, moriverDataEntry.GearType, moriverDataEntry.Temp, moriverDataEntry.Turbidity, moriverDataEntry.Conductivity, moriverDataEntry.Do,
		moriverDataEntry.Distance, moriverDataEntry.Width, moriverDataEntry.Netrivermile, moriverDataEntry.Structurenumber, moriverDataEntry.Usgs, moriverDataEntry.Riverstage, moriverDataEntry.Discharge,
		moriverDataEntry.U1, moriverDataEntry.U2, moriverDataEntry.U3, moriverDataEntry.U4, moriverDataEntry.U5, moriverDataEntry.U6, moriverDataEntry.U7, moriverDataEntry.Macro, moriverDataEntry.Meso, moriverDataEntry.Habitatrn, moriverDataEntry.Qc,
		moriverDataEntry.MicroStructure, moriverDataEntry.StructureFlow, moriverDataEntry.StructureMod, moriverDataEntry.SetSite1, moriverDataEntry.SetSite2, moriverDataEntry.SetSite3,
		moriverDataEntry.StartTime, moriverDataEntry.StartLatitude, moriverDataEntry.StartLongitude, moriverDataEntry.StopTime, moriverDataEntry.StopLatitude, moriverDataEntry.StopLongitude,
		moriverDataEntry.Depth1, moriverDataEntry.Velocitybot1, moriverDataEntry.Velocity08_1, moriverDataEntry.Velocity02or06_1,
		moriverDataEntry.Depth2, moriverDataEntry.Velocitybot2, moriverDataEntry.Velocity08_2, moriverDataEntry.Velocity02or06_2,
		moriverDataEntry.Depth3, moriverDataEntry.Velocitybot3, moriverDataEntry.Velocity08_3, moriverDataEntry.Velocity02or06_3,
		moriverDataEntry.Watervel, moriverDataEntry.Cobble, moriverDataEntry.Organic, moriverDataEntry.Silt, moriverDataEntry.Sand, moriverDataEntry.Gravel,
		moriverDataEntry.Comments, moriverDataEntry.Complete, moriverDataEntry.Checkby, moriverDataEntry.NoTurbidity, moriverDataEntry.NoVelocity, moriverDataEntry.EditInitials, moriverDataEntry.MrFid, moriverDataEntry.SiteID, moriverDataEntry.FieldOffice,
		moriverDataEntry.LastEditComment, moriverDataEntry.LastUpdated, moriverDataEntry.UploadedBy, moriverDataEntry.Bend, moriverDataEntry.BendRn, moriverDataEntry.BendRiverMile, moriverDataEntry.MrID)
	return err
}

var moriverDataEntriesByFidSql = `select mr_fid,mr_id,site_id,FIELDOFFICE,PROJECT,SEGMENT,SEASON,setdate, subsample, subsamplepass, 
subsamplen, recorder, gear, GEAR_TYPE, temp, turbidity, conductivity, do, distance, width, netrivermile, structurenumber, usgs, riverstage, discharge,
u1, u2, u3, u4, u5, u6, u7, MACRO, MESO, habitatrn, qc, micro_structure, structure_flow, structure_mod, set_site_1, set_site_2, set_site_3,
starttime, startlatitude, startlongitude, stoptime, stoplatitude, stoplongitude, depth1, velocitybot1, velocity08_1, velocity02or06_1,
depth2, velocitybot2, velocity08_2, velocity02or06_2,depth3, velocitybot3, velocity08_3, velocity02or06_3, watervel, cobble, ORGANIC, silt, sand,
gravel, comments, complete, checkby, no_turbidity, no_velocity, edit_initials,last_edit_comment, uploaded_by from ds_moriver 
where mr_id = :1`

var moriverDataEntriesCountByFidSql = `SELECT count(*) FROM ds_moriver where mr_id = :1`

var moriverDataEntriesByFfidSql = `select mr_fid,mr_id,site_id,FIELDOFFICE,PROJECT,SEGMENT,SEASON,setdate, subsample, subsamplepass, 
subsamplen, recorder, gear, GEAR_TYPE, temp, turbidity, conductivity, do, distance, width, netrivermile, structurenumber, usgs, riverstage, discharge,
u1, u2, u3, u4, u5, u6, u7, MACRO, MESO, habitatrn, qc, micro_structure, structure_flow, structure_mod, set_site_1, set_site_2, set_site_3,
starttime, startlatitude, startlongitude, stoptime, stoplatitude, stoplongitude, depth1, velocitybot1, velocity08_1, velocity02or06_1,
depth2, velocitybot2, velocity08_2, velocity02or06_2,depth3, velocitybot3, velocity08_3, velocity02or06_3, watervel, cobble, ORGANIC, silt, sand,
gravel, comments, complete, checkby, no_turbidity, no_velocity, edit_initials,last_edit_comment, uploaded_by from ds_moriver 
where mr_fid = :1`

var moriverDataEntriesCountByFfidSql = `SELECT count(*) FROM ds_moriver where mr_fid = :1`

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

	countrows, err := countQuery.Query(id)
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
		queryParams.OrderBy = "mr_id desc"
	}
	moriverDataEntriesSqlWithSearch := query + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(moriverDataEntriesSqlWithSearch)
	if err != nil {
		return moriverDataEntryWithCount, err
	}

	rows, err := dbQuery.Query(id)
	if err != nil {
		return moriverDataEntryWithCount, err
	}
	defer rows.Close()

	for rows.Next() {
		moriverDataEntry := models.UploadMoriver{}
		err = rows.Scan(&moriverDataEntry.MrFid, &moriverDataEntry.MrID, &moriverDataEntry.SiteID, &moriverDataEntry.FieldOffice,
			&moriverDataEntry.Project, &moriverDataEntry.Segment, &moriverDataEntry.Season, &moriverDataEntry.SetDate, &moriverDataEntry.Subsample, &moriverDataEntry.Subsamplepass,
			&moriverDataEntry.Subsamplen, &moriverDataEntry.Recorder, &moriverDataEntry.Gear, &moriverDataEntry.GearType, &moriverDataEntry.Temp, &moriverDataEntry.Turbidity, &moriverDataEntry.Conductivity, &moriverDataEntry.Do,
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

var insertSupplementalDataSql = `insert into ds_supplemental(s_id, f_id, f_fid, mr_id,TAGNUMBER, PITRN,SCUTELOC, SCUTENUM, SCUTELOC2, SCUTENUM2,ELHV, ELCOLOR, ERHV, ERCOLOR, CWTYN, DANGLER, genetic, genetics_vial_number,
BROODSTOCK, HATCH_WILD, species_id,head, snouttomouth, inter, mouthwidth, m_ib,l_ob, l_ib, r_ib,r_ob, anal, dorsal, status, HATCHERY_ORIGIN,SEX, stage, recapture, photo,genetic_needs, other_tag_info,comments,
edit_initials,last_edit_comment, last_updated, uploaded_by) 
values (supp_seq.nextval,:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23,:24,:25,:26,:27,:28,:29,:30,:31,:32,:33,:34,:35,:36,:37,:38,:39,:40,:41,:42,:43,:44)`

func (s *PallidSturgeonStore) SaveSupplementalDataEntry(supplementalDataEntry models.UploadSupplemental) error {
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
		supplementalDataEntry.UploadedBy)
	return err
}

var updateSupplementalDataSql = `UPDATE ds_supplemental SET 
f_fid = :2,
mr_id = :3,
TAGNUMBER = :4, 
PITRN = :5, 
SCUTELOC = :6, 
SCUTENUM = :7, 
SCUTELOC2  = :8, 
SCUTENUM2 = :9, 
ELHV = :10, 
ELCOLOR = :11, 
ERHV = :12, 
ERCOLOR = :13, 
CWTYN  = :14, 
DANGLER = :15, 
genetic = :16, 
genetics_vial_number = :17,
BROODSTOCK = :18, 
HATCH_WILD = :19, 
head = :20, 
snouttomouth = :21, 
inter = :22, 
mouthwidth = :23, 
m_ib = :24, 
l_ob = :25, 
l_ib = :26, 
r_ib = :27, 
r_ob = :28, 
anal = :29, 
dorsal = :30, 
status = :31, 
HATCHERY_ORIGIN = :32, 
SEX = :33, 
stage = :34,  
recapture = :35, 
photo = :36, 
genetic_needs = :37, 
other_tag_info = :38, 
comments = :39, 
f_id = :40,
edit_initials = :41,
last_edit_comment = :42,	
last_updated = :43, 
uploaded_by = :44,
species_id = :45
WHERE s_id = :1`

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
		supplementalDataEntry.SpeciesId,
		supplementalDataEntry.Sid)
	return err
}

var supplementalDataEntriesSql = `select su.s_id, su.f_id, su.f_fid, su.mr_id, si.site_id, su.tagnumber, su.pitrn, su.scuteloc, su.scutenum, su.scuteloc2, su.scutenum2, su.elhv, su.elcolor, su.erhv, su.ercolor, su.cwtyn, 
su.dangler, su.genetic, su.genetics_vial_number, su.broodstock, su.hatch_wild, su.species_id, su.head, su.snouttomouth, su.inter, su.mouthwidth, su.m_ib, su.l_ob, su.l_ib, su.r_ib, su.r_ob, su.anal, su.dorsal, su.status, 
su.hatchery_origin, su.sex, su.stage, su.recapture, su.photo, su.genetic_needs, su.other_tag_info, su.comments, su.edit_initials, su.last_edit_comment, su.uploaded_by from ds_supplemental su
inner join ds_moriver mo on su.mr_id = mo.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :1 != 'ZZ' THEN si.fieldoffice ELSE :2 END) = :3`

var supplementalDataEntriesCountBySql = `SELECT count(*) from ds_supplemental su
inner join ds_moriver mo on su.mr_id = mo.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :1 != 'ZZ' THEN si.fieldoffice ELSE :2 END) = :3`

var supplementalDataEntriesByFidSql = `select su.s_id, su.f_id, su.f_fid, su.mr_id, si.site_id, su.tagnumber, su.pitrn, su.scuteloc, su.scutenum, su.scuteloc2, su.scutenum2, su.elhv, su.elcolor, su.erhv, su.ercolor, su.cwtyn, 
su.dangler, su.genetic, su.genetics_vial_number, su.broodstock, su.hatch_wild, su.species_id, su.head, su.snouttomouth, su.inter, su.mouthwidth, su.m_ib, su.l_ob, su.l_ib, su.r_ib, su.r_ob, su.anal, su.dorsal, su.status, 
su.hatchery_origin, su.sex, su.stage, su.recapture, su.photo, su.genetic_needs, su.other_tag_info, su.comments, su.edit_initials, su.last_edit_comment, su.uploaded_by from ds_supplemental su
inner join ds_moriver mo on su.mr_id = mo.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4
and su.f_id = :1`

var supplementalDataEntriesCountByFidSql = `SELECT count(*) from ds_supplemental su
inner join ds_moriver mo on su.mr_id = mo.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4
and su.f_id = :1`

var supplementalDataEntriesByFfidSql = `select su.s_id, su.f_id, su.f_fid, su.mr_id, si.site_id, su.tagnumber, su.pitrn, su.scuteloc, su.scutenum, su.scuteloc2, su.scutenum2, su.elhv, su.elcolor, su.erhv, su.ercolor, su.cwtyn, 
su.dangler, su.genetic, su.genetics_vial_number, su.broodstock, su.hatch_wild, su.species_id, su.head, su.snouttomouth, su.inter, su.mouthwidth, su.m_ib, su.l_ob, su.l_ib, su.r_ib, su.r_ob, su.anal, su.dorsal, su.status, 
su.hatchery_origin, su.sex, su.stage, su.recapture, su.photo, su.genetic_needs, su.other_tag_info, su.comments, su.edit_initials, su.last_edit_comment, su.uploaded_by from ds_supplemental su
inner join ds_moriver mo on su.mr_id = mo.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4
and su.f_fid = :1`

var supplementalDataEntriesCountByFfidSql = `SELECT count(*) from ds_supplemental su
inner join ds_moriver mo on su.mr_id = mo.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4
and su.f_fid = :1`

var supplementalDataEntriesByGeneticsVialSql = `select su.s_id, su.f_id, su.f_fid, su.mr_id, si.site_id, su.tagnumber, su.pitrn, su.scuteloc, su.scutenum, su.scuteloc2, su.scutenum2, su.elhv, su.elcolor, su.erhv, su.ercolor, su.cwtyn, 
su.dangler, su.genetic, su.genetics_vial_number, su.broodstock, su.hatch_wild, su.species_id, su.head, su.snouttomouth, su.inter, su.mouthwidth, su.m_ib, su.l_ob, su.l_ib, su.r_ib, su.r_ob, su.anal, su.dorsal, su.status, 
su.hatchery_origin, su.sex, su.stage, su.recapture, su.photo, su.genetic_needs, su.other_tag_info, su.comments, su.edit_initials, su.last_edit_comment, su.uploaded_by from ds_supplemental su
inner join ds_moriver mo on su.mr_id = mo.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4
and su.genetics_vial_number = :1`

var supplementalDataEntriesCountByGeneticsVialSql = `SELECT count(*) from ds_supplemental su
inner join ds_moriver mo on su.mr_id = mo.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4
and su.genetics_vial_number = :1`

var supplementalDataEntriesByGeneticsPitTagSql = `select su.s_id, su.f_id, su.f_fid, su.mr_id, si.site_id, su.tagnumber, su.pitrn, su.scuteloc, su.scutenum, su.scuteloc2, su.scutenum2, su.elhv, su.elcolor, su.erhv, su.ercolor, su.cwtyn, 
su.dangler, su.genetic, su.genetics_vial_number, su.broodstock, su.hatch_wild, su.species_id, su.head, su.snouttomouth, su.inter, su.mouthwidth, su.m_ib, su.l_ob, su.l_ib, su.r_ib, su.r_ob, su.anal, su.dorsal, su.status, 
su.hatchery_origin, su.sex, su.stage, su.recapture, su.photo, su.genetic_needs, su.other_tag_info, su.comments, su.edit_initials, su.last_edit_comment, su.uploaded_by from ds_supplemental su
inner join ds_moriver mo on su.mr_id = mo.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4
and su.TAGNUMBER = :1`

var supplementalDataEntriesCountByPitTagSql = `SELECT count(*) from ds_supplemental su
inner join ds_moriver mo on su.mr_id = mo.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4
and su.TAGNUMBER = :1`

var supplementalDataEntriesByMrIdSql = `select su.s_id, su.f_id, su.f_fid, su.mr_id, si.site_id, su.tagnumber, su.pitrn, su.scuteloc, su.scutenum, su.scuteloc2, su.scutenum2, su.elhv, su.elcolor, su.erhv, su.ercolor, su.cwtyn, 
su.dangler, su.genetic, su.genetics_vial_number, su.broodstock, su.hatch_wild, su.species_id, su.head, su.snouttomouth, su.inter, su.mouthwidth, su.m_ib, su.l_ob, su.l_ib, su.r_ib, su.r_ob, su.anal, su.dorsal, su.status, 
su.hatchery_origin, su.sex, su.stage, su.recapture, su.photo, su.genetic_needs, su.other_tag_info, su.comments, su.edit_initials, su.last_edit_comment, su.uploaded_by from ds_supplemental su
inner join ds_moriver mo on su.mr_id = mo.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4
and su.mr_id = :1`

var supplementalDataEntriesCountByMrIdSql = `SELECT count(*) from ds_supplemental su
inner join ds_moriver mo on su.mr_id = mo.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4
and su.mr_id = :1`

var supplementalDataEntriesBySidSql = `select su.s_id, su.f_id, su.f_fid, su.mr_id, si.site_id, su.tagnumber, su.pitrn, su.scuteloc, su.scutenum, su.scuteloc2, su.scutenum2, su.elhv, su.elcolor, su.erhv, su.ercolor, su.cwtyn, 
su.dangler, su.genetic, su.genetics_vial_number, su.broodstock, su.hatch_wild, su.species_id, su.head, su.snouttomouth, su.inter, su.mouthwidth, su.m_ib, su.l_ob, su.l_ib, su.r_ib, su.r_ob, su.anal, su.dorsal, su.status, 
su.hatchery_origin, su.sex, su.stage, su.recapture, su.photo, su.genetic_needs, su.other_tag_info, su.comments, su.edit_initials, su.last_edit_comment, su.uploaded_by from ds_supplemental su
inner join ds_moriver mo on su.mr_id = mo.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4
and su.s_id = :1`

var supplementalDataEntriesCountBySidSql = `SELECT count(*) from ds_supplemental su
inner join ds_moriver mo on su.mr_id = mo.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4
and su.s_id = :1`

func (s *PallidSturgeonStore) GetSupplementalDataEntries(tableId string, fieldId string, geneticsVial string, pitTag string, mrId string, fId string, officeCode string, queryParams models.SearchParams) (models.SupplementalDataEntryWithCount, error) {
	supplementalDataEntryWithCount := models.SupplementalDataEntryWithCount{}
	query := ""
	queryWithCount := ""
	id := ""

	if tableId != "" {
		query = supplementalDataEntriesBySidSql
		queryWithCount = supplementalDataEntriesCountBySidSql
		id = tableId
	}

	if fId != "" {
		query = supplementalDataEntriesByFidSql
		queryWithCount = supplementalDataEntriesCountByFidSql
		id = fId
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

	if tableId == "" && fId == "" && fieldId == "" && geneticsVial == "" && pitTag == "" && mrId == "" {
		query = supplementalDataEntriesSql
		queryWithCount = supplementalDataEntriesCountBySql
	}

	countQuery, err := s.db.Prepare(queryWithCount)
	if err != nil {
		return supplementalDataEntryWithCount, err
	}

	var countrows *sql.Rows
	if id == "" {
		countrows, err = countQuery.Query(officeCode, officeCode, officeCode)
		if err != nil {
			return supplementalDataEntryWithCount, err
		}
	} else {
		countrows, err = countQuery.Query(officeCode, officeCode, officeCode, id)
		if err != nil {
			return supplementalDataEntryWithCount, err
		}
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
		queryParams.OrderBy = "s_id"
	}
	supplementalDataEntriesSqlWithSearch := query + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(supplementalDataEntriesSqlWithSearch)
	if err != nil {
		return supplementalDataEntryWithCount, err
	}

	var rows *sql.Rows
	if id == "" {
		rows, err = dbQuery.Query(officeCode, officeCode, officeCode)
		if err != nil {
			return supplementalDataEntryWithCount, err
		}
	} else {
		rows, err = dbQuery.Query(officeCode, officeCode, officeCode, id)
		if err != nil {
			return supplementalDataEntryWithCount, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		supplementalDataEntry := models.UploadSupplemental{}
		err = rows.Scan(
			&supplementalDataEntry.Sid,
			&supplementalDataEntry.Fid,
			&supplementalDataEntry.FFid,
			&supplementalDataEntry.MrId,
			&supplementalDataEntry.SiteID,
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

var searchDataEntriesSql = `select SE_FID, SE_ID, CHECKBY, conductivity, EDIT_INITIALS, LAST_EDIT_COMMENT, LAST_UPDATED, RECORDER, SEARCH_DATE, search_day, 
SEARCH_TYPE_CODE, SITE_ID, START_LATITUDE, START_LONGITUDE, START_TIME, STOP_LATITUDE, STOP_LONGITUDE, STOP_TIME, temp, UPLOADED_BY, UPLOAD_FILENAME,
UPLOAD_SESSION_ID, ds_id from ds_search`

var searchDataEntriesCountSql = `select count(*) from ds_search`

var searchDataEntriesBySeIdSql = `select SE_FID, SE_ID, CHECKBY, conductivity, EDIT_INITIALS, LAST_EDIT_COMMENT, LAST_UPDATED, RECORDER, SEARCH_DATE, search_day, 
SEARCH_TYPE_CODE, SITE_ID, START_LATITUDE, START_LONGITUDE, START_TIME, STOP_LATITUDE, STOP_LONGITUDE, STOP_TIME, temp, UPLOADED_BY, UPLOAD_FILENAME,
UPLOAD_SESSION_ID, ds_id from ds_search where se_id = :1`

var searchDataEntriesCountBySeIdSql = `select count(*) from ds_search where se_id = :1`

var searchDataEntriesBySiteIdSql = `select SE_FID, SE_ID, CHECKBY, conductivity, EDIT_INITIALS, LAST_EDIT_COMMENT, LAST_UPDATED, RECORDER, SEARCH_DATE, search_day, 
SEARCH_TYPE_CODE, SITE_ID, START_LATITUDE, START_LONGITUDE, START_TIME, STOP_LATITUDE, STOP_LONGITUDE, STOP_TIME, temp, UPLOADED_BY, UPLOAD_FILENAME,
UPLOAD_SESSION_ID, ds_id from ds_search where site_id = :1`

var searchDataEntriesCountBySiteIdSql = `select count(*) from ds_search where site_id = :1`

func (s *PallidSturgeonStore) GetSearchDataEntries(tableId string, siteId string, queryParams models.SearchParams) (models.SearchDataEntryWithCount, error) {
	searchDataEntryWithCount := models.SearchDataEntryWithCount{}
	query := ""
	queryWithCount := ""
	id := ""

	if tableId != "" {
		query = searchDataEntriesBySeIdSql
		queryWithCount = searchDataEntriesCountBySeIdSql
		id = tableId
	}

	if siteId != "" {
		query = searchDataEntriesBySiteIdSql
		queryWithCount = searchDataEntriesCountBySiteIdSql
		id = siteId
	}

	if tableId == "" && siteId == "" {
		query = searchDataEntriesSql
		queryWithCount = searchDataEntriesCountSql
	}

	countQuery, err := s.db.Prepare(queryWithCount)
	if err != nil {
		return searchDataEntryWithCount, err
	}

	var countrows *sql.Rows
	if id == "" {
		countrows, err = countQuery.Query()
		if err != nil {
			return searchDataEntryWithCount, err
		}
	} else {
		countrows, err = countQuery.Query(id)
		if err != nil {
			return searchDataEntryWithCount, err
		}
	}
	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&searchDataEntryWithCount.TotalCount)
		if err != nil {
			return searchDataEntryWithCount, err
		}
	}

	searchEntries := []models.UploadSearch{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "se_id desc"
	}
	searchDataEntriesSqlWithSearch := query + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(searchDataEntriesSqlWithSearch)
	if err != nil {
		return searchDataEntryWithCount, err
	}

	var rows *sql.Rows
	if id == "" {
		rows, err = dbQuery.Query()
		if err != nil {
			return searchDataEntryWithCount, err
		}
	} else {
		rows, err = dbQuery.Query(id)
		if err != nil {
			return searchDataEntryWithCount, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		searchDataEntry := models.UploadSearch{}
		err = rows.Scan(&searchDataEntry.SeFid, &searchDataEntry.SeId, &searchDataEntry.Checkby, &searchDataEntry.Conductivity, &searchDataEntry.EditInitials, &searchDataEntry.LastEditComment, &searchDataEntry.LastUpdated,
			&searchDataEntry.Recorder, &searchDataEntry.SearchDate, &searchDataEntry.SearchDay, &searchDataEntry.SearchTypeCode, &searchDataEntry.SiteId, &searchDataEntry.StartLatitude, &searchDataEntry.StartLongitude,
			&searchDataEntry.StartTime, &searchDataEntry.StopLatitude, &searchDataEntry.StopLongitude, &searchDataEntry.StopTime, &searchDataEntry.Temp, &searchDataEntry.UploadedBy, &searchDataEntry.UploadFilename,
			&searchDataEntry.UploadSessionId, &searchDataEntry.DsId)
		if err != nil {
			return searchDataEntryWithCount, err
		}
		searchEntries = append(searchEntries, searchDataEntry)
	}

	searchDataEntryWithCount.Items = searchEntries

	return searchDataEntryWithCount, err
}

var insertSearchDataSql = `insert into ds_search (se_id, SE_FID, CHECKBY, conductivity, EDIT_INITIALS, LAST_EDIT_COMMENT, LAST_UPDATED, RECORDER, SEARCH_DATE,
SEARCH_TYPE_CODE, SITE_ID, START_LATITUDE, START_LONGITUDE, START_TIME, STOP_LATITUDE, STOP_LONGITUDE, STOP_TIME, temp, UPLOADED_BY, UPLOAD_FILENAME,
UPLOAD_SESSION_ID, ds_id) values (search_seq.nextval,:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21)`

func (s *PallidSturgeonStore) SaveSearchDataEntry(searchDataEntry models.UploadSearch) error {
	_, err := s.db.Exec(insertSearchDataSql, searchDataEntry.SeFid, searchDataEntry.Checkby, searchDataEntry.Conductivity, searchDataEntry.EditInitials, searchDataEntry.LastEditComment, searchDataEntry.LastUpdated, searchDataEntry.Recorder,
		searchDataEntry.SearchDate, searchDataEntry.SearchTypeCode, searchDataEntry.SiteId, searchDataEntry.StartLatitude, searchDataEntry.StartLongitude, searchDataEntry.StartTime, searchDataEntry.StopLatitude,
		searchDataEntry.StopLongitude, searchDataEntry.StopTime, searchDataEntry.Temp, searchDataEntry.UploadedBy, searchDataEntry.UploadFilename, searchDataEntry.UploadSessionId, searchDataEntry.DsId)
	return err
}

var updateSearchDataSql = `UPDATE ds_search SET 
SE_FID = :2,
CHECKBY = :3,
CONDUCTIVITY = :4,
EDIT_INITIALS = :5,
LAST_EDIT_COMMENT = :6,
LAST_UPDATED = :7,
RECORDER = :8,
SEARCH_DATE = :9,
SEARCH_DAY = :10,
SEARCH_TYPE_CODE = :11,
SITE_ID = :12,
START_LATITUDE = :13,
START_LONGITUDE = :14,
START_TIME = :15,
STOP_LATITUDE = :16,
STOP_LONGITUDE = :17,
STOP_TIME = :18,
TEMP = :19,
UPLOADED_BY = :20,
UPLOAD_FILENAME = :21,
UPLOAD_SESSION_ID = :22
WHERE SE_ID = :1`

func (s *PallidSturgeonStore) UpdateSearchDataEntry(searchDataEntry models.UploadSearch) error {
	_, err := s.db.Exec(updateSearchDataSql, searchDataEntry.SeFid, searchDataEntry.Checkby, searchDataEntry.Conductivity, searchDataEntry.EditInitials, searchDataEntry.LastEditComment, searchDataEntry.LastUpdated, searchDataEntry.Recorder,
		searchDataEntry.SearchDate, searchDataEntry.SearchDay, searchDataEntry.SearchTypeCode, searchDataEntry.SiteId, searchDataEntry.StartLatitude, searchDataEntry.StartLongitude, searchDataEntry.StartTime, searchDataEntry.StopLatitude,
		searchDataEntry.StopLongitude, searchDataEntry.StopTime, searchDataEntry.Temp, searchDataEntry.UploadedBy, searchDataEntry.UploadFilename, searchDataEntry.UploadSessionId, searchDataEntry.SeId)
	return err
}

var telemetryDataEntriesSql = `select te.bend,te.CAPTURE_LATITUDE,te.CAPTURE_LONGITUDE,te.CAPTURE_TIME,te.CHECKBY,te.COMMENTS,te.conductivity,te.depth,te.EDIT_INITIALS,
te.FREQUENCY_ID_CODE,te.gravel,te.LAST_EDIT_COMMENT,te.LAST_UPDATED,te.MACRO_ID,te.MESO_ID,te.position_confidence,te.RADIO_TAG_NUM,
te.sand,te.SE_FID,te.SE_ID,te.silt,te.temp,te.turbidity,te.T_FID,te.T_ID,te.UPLOADED_BY,te.UPLOAD_FILENAME,
te.UPLOAD_SESSION_ID, si.site_id from ds_telemetry_fish te
inner join ds_search se on te.se_id = se.se_id
inner join ds_sites si on si.site_id = se.site_id
where (CASE when :1 != 'ZZ' THEN si.fieldoffice ELSE :2 END) = :3`

var telemetryDataEntriesCountSql = `select count(*) from ds_telemetry_fish te
inner join ds_search se on te.se_id = se.se_id
inner join ds_sites si on si.site_id = se.site_id
where (CASE when :1 != 'ZZ' THEN si.fieldoffice ELSE :2 END) = :3`

var telemetryDataEntriesBySeIdSql = `select te.bend,te.CAPTURE_LATITUDE,te.CAPTURE_LONGITUDE,te.CAPTURE_TIME,te.CHECKBY,te.COMMENTS,te.conductivity,te.depth,te.EDIT_INITIALS,
te.FREQUENCY_ID_CODE,te.gravel,te.LAST_EDIT_COMMENT,te.LAST_UPDATED,te.MACRO_ID,te.MESO_ID,te.position_confidence,te.RADIO_TAG_NUM,
te.sand,te.SE_FID,te.SE_ID,te.silt,te.temp,te.turbidity,te.T_FID,te.T_ID,te.UPLOADED_BY,te.UPLOAD_FILENAME,
te.UPLOAD_SESSION_ID, si.site_id from ds_telemetry_fish te
inner join ds_search se on te.se_id = se.se_id
inner join ds_sites si on si.site_id = se.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4 
and te.se_id = :1`

var telemetryDataEntriesCountBySeIdSql = `select count(*) from ds_telemetry_fish te
inner join ds_search se on te.se_id = se.se_id
inner join ds_sites si on si.site_id = se.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4 
and te.se_id = :1`

var telemetryDataEntriesByTidSql = `select te.bend,te.CAPTURE_LATITUDE,te.CAPTURE_LONGITUDE,te.CAPTURE_TIME,te.CHECKBY,te.COMMENTS,te.conductivity,te.depth,te.EDIT_INITIALS,
te.FREQUENCY_ID_CODE,te.gravel,te.LAST_EDIT_COMMENT,te.LAST_UPDATED,te.MACRO_ID,te.MESO_ID,te.position_confidence,te.RADIO_TAG_NUM,
te.sand,te.SE_FID,te.SE_ID,te.silt,te.temp,te.turbidity,te.T_FID,te.T_ID,te.UPLOADED_BY,te.UPLOAD_FILENAME,
te.UPLOAD_SESSION_ID, si.site_id from ds_telemetry_fish te
inner join ds_search se on te.se_id = se.se_id
inner join ds_sites si on si.site_id = se.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4 
and te.t_id = :1`

var telemetryDataEntriesCountByTidSql = `select count(*) from ds_telemetry_fish te
inner join ds_search se on te.se_id = se.se_id
inner join ds_sites si on si.site_id = se.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4 
and te.t_id = :1`

func (s *PallidSturgeonStore) GetTelemetryDataEntries(tableId string, seId string, officeCode string, queryParams models.SearchParams) (models.TelemetryDataEntryWithCount, error) {
	telemetryDataEntryWithCount := models.TelemetryDataEntryWithCount{}
	query := ""
	queryWithCount := ""
	id := ""

	if tableId != "" {
		query = telemetryDataEntriesByTidSql
		queryWithCount = telemetryDataEntriesCountByTidSql
		id = tableId
	}

	if seId != "" {
		query = telemetryDataEntriesBySeIdSql
		queryWithCount = telemetryDataEntriesCountBySeIdSql
		id = seId
	}

	if tableId == "" && seId == "" {
		query = telemetryDataEntriesSql
		queryWithCount = telemetryDataEntriesCountSql
	}

	countQuery, err := s.db.Prepare(queryWithCount)
	if err != nil {
		return telemetryDataEntryWithCount, err
	}

	var countrows *sql.Rows
	if id == "" {
		countrows, err = countQuery.Query(officeCode, officeCode, officeCode)
		if err != nil {
			return telemetryDataEntryWithCount, err
		}
	} else {
		countrows, err = countQuery.Query(officeCode, officeCode, officeCode, id)
		if err != nil {
			return telemetryDataEntryWithCount, err
		}
	}
	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&telemetryDataEntryWithCount.TotalCount)
		if err != nil {
			return telemetryDataEntryWithCount, err
		}
	}

	telemetryEntries := []models.UploadTelemetry{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "t_id desc"
	}
	telemetryDataEntriesSqlWithSearch := query + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(telemetryDataEntriesSqlWithSearch)
	if err != nil {
		return telemetryDataEntryWithCount, err
	}

	var rows *sql.Rows
	if id == "" {
		rows, err = dbQuery.Query(officeCode, officeCode, officeCode)
		if err != nil {
			return telemetryDataEntryWithCount, err
		}
	} else {
		rows, err = dbQuery.Query(officeCode, officeCode, officeCode, id)
		if err != nil {
			return telemetryDataEntryWithCount, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		telemetryDataEntry := models.UploadTelemetry{}
		err = rows.Scan(&telemetryDataEntry.Bend, &telemetryDataEntry.CaptureLatitude, &telemetryDataEntry.CaptureLongitude, &telemetryDataEntry.CaptureTime, &telemetryDataEntry.Checkby, &telemetryDataEntry.Comments, &telemetryDataEntry.Conductivity, &telemetryDataEntry.Depth,
			&telemetryDataEntry.EditInitials, &telemetryDataEntry.FrequencyIdCode, &telemetryDataEntry.Gravel, &telemetryDataEntry.LastEditComment, &telemetryDataEntry.LastUpdated, &telemetryDataEntry.MacroId, &telemetryDataEntry.MesoId, &telemetryDataEntry.PositionConfidence,
			&telemetryDataEntry.RadioTagNum, &telemetryDataEntry.Sand, &telemetryDataEntry.SeFid, &telemetryDataEntry.SeId, &telemetryDataEntry.Silt, &telemetryDataEntry.Temp, &telemetryDataEntry.Turbidity, &telemetryDataEntry.TFid, &telemetryDataEntry.TId, &telemetryDataEntry.UploadedBy,
			&telemetryDataEntry.UploadFilename, &telemetryDataEntry.UploadSessionId, &telemetryDataEntry.SiteId)
		if err != nil {
			return telemetryDataEntryWithCount, err
		}
		telemetryEntries = append(telemetryEntries, telemetryDataEntry)
	}

	telemetryDataEntryWithCount.Items = telemetryEntries

	return telemetryDataEntryWithCount, err
}

var insertTelemetryDataSql = `insert into ds_telemetry_fish (t_id, BEND,CAPTURE_LATITUDE,CAPTURE_LONGITUDE,CAPTURE_TIME,CHECKBY,COMMENTS,conductivity,depth,EDIT_INITIALS,FREQUENCY_ID_CODE,gravel,LAST_EDIT_COMMENT,
	LAST_UPDATED,MACRO_ID,MESO_ID,POSITION_CONFIDENCE,RADIO_TAG_NUM,sand,SE_FID,SE_ID,silt,temp,turbidity,T_FID,UPLOADED_BY,UPLOAD_FILENAME,UPLOAD_SESSION_ID) 
	values (telemetry_id_seq.nextval, :1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23,:24,:25,:26,:27)`

func (s *PallidSturgeonStore) SaveTelemetryDataEntry(telemetryDataEntry models.UploadTelemetry) error {
	_, err := s.db.Exec(insertTelemetryDataSql, telemetryDataEntry.Bend, telemetryDataEntry.CaptureLatitude, telemetryDataEntry.CaptureLongitude, telemetryDataEntry.CaptureTime, telemetryDataEntry.Checkby, telemetryDataEntry.Comments, telemetryDataEntry.Conductivity, telemetryDataEntry.Depth,
		telemetryDataEntry.EditInitials, telemetryDataEntry.FrequencyIdCode, telemetryDataEntry.Gravel, telemetryDataEntry.LastEditComment, telemetryDataEntry.LastUpdated, telemetryDataEntry.MacroId, telemetryDataEntry.MesoId, telemetryDataEntry.PositionConfidence, telemetryDataEntry.RadioTagNum,
		telemetryDataEntry.Sand, telemetryDataEntry.SeFid, telemetryDataEntry.SeId, telemetryDataEntry.Silt, telemetryDataEntry.Temp, telemetryDataEntry.Turbidity, telemetryDataEntry.TFid, telemetryDataEntry.UploadedBy, telemetryDataEntry.UploadFilename, telemetryDataEntry.UploadSessionId)
	return err
}

var updateTelemetryDataSql = `UPDATE ds_telemetry_fish SET 
BEND = :2,
CAPTURE_LATITUDE = :3,
CAPTURE_LONGITUDE = :4,
CAPTURE_TIME = :5,
CHECKBY = :6,
COMMENTS = :7,
CONDUCTIVITY = :8,
DEPTH = :9,
EDIT_INITIALS = :10,
FREQUENCY_ID_CODE = :11,
GRAVEL = :12,
LAST_EDIT_COMMENT = :13,
LAST_UPDATED = :14,
MACRO_ID = :15,
MESO_ID = :16,
POSITION_CONFIDENCE = :17,
RADIO_TAG_NUM = :18,
SAND = :19,
SE_FID = :20,
SE_ID = :21,
SILT = :22,
TEMP = :23,
TURBIDITY = :24,
T_FID = :25,
UPLOADED_BY = :26,
UPLOAD_FILENAME = :27,
UPLOAD_SESSION_ID = :28
WHERE T_ID = :1`

func (s *PallidSturgeonStore) UpdateTelemetryDataEntry(telemetryDataEntry models.UploadTelemetry) error {
	_, err := s.db.Exec(updateTelemetryDataSql, telemetryDataEntry.Bend, telemetryDataEntry.CaptureLatitude, telemetryDataEntry.CaptureLongitude, telemetryDataEntry.CaptureTime, telemetryDataEntry.Checkby, telemetryDataEntry.Comments, telemetryDataEntry.Conductivity, telemetryDataEntry.Depth,
		telemetryDataEntry.EditInitials, telemetryDataEntry.FrequencyIdCode, telemetryDataEntry.Gravel, telemetryDataEntry.LastEditComment, telemetryDataEntry.LastUpdated, telemetryDataEntry.MacroId, telemetryDataEntry.MesoId, telemetryDataEntry.PositionConfidence, telemetryDataEntry.RadioTagNum,
		telemetryDataEntry.Sand, telemetryDataEntry.SeFid, telemetryDataEntry.SeId, telemetryDataEntry.Silt, telemetryDataEntry.Temp, telemetryDataEntry.Turbidity, telemetryDataEntry.TFid, telemetryDataEntry.UploadedBy, telemetryDataEntry.UploadFilename, telemetryDataEntry.UploadSessionId, telemetryDataEntry.TId)
	return err
}

var procedureDataEntriesSql = `select pr.ID, pr.F_ID, pr.F_FID, si.site_id, pr.PURPOSE_CODE, pr.PROCEDURE_DATE, pr.PROCEDURE_START_TIME, pr.PROCEDURE_END_TIME, pr.PROCEDURE_BY, pr.ANTIBIOTIC_INJECTION_IND, pr.PHOTO_DORSAL_IND, pr.PHOTO_VENTRAL_IND, 
pr.PHOTO_LEFT_IND, pr.OLD_RADIO_TAG_NUM, pr.OLD_FREQUENCY_ID, pr.DST_SERIAL_NUM, pr.DST_START_TIME, pr.DST_REIMPLANT_IND, pr.NEW_RADIO_TAG_NUM, pr.NEW_FREQUENCY_ID, pr.SEX_CODE, pr.COMMENTS, pr.FISH_HEALTH_COMMENTS, pr.SPAWN_CODE, pr.EVAL_LOCATION_CODE, 
pr.BLOOD_SAMPLE_IND, pr.EGG_SAMPLE_IND, pr.VISUAL_REPRO_STATUS_CODE, pr.ULTRASOUND_REPRO_STATUS_CODE, pr.ULTRASOUND_GONAD_LENGTH, pr.GONAD_CONDITION, pr.EXPECTED_SPAWN_YEAR, pr.LAST_UPDATED, pr.UPLOAD_SESSION_ID, pr.UPLOADED_BY, pr.UPLOAD_FILENAME, 
pr.CHECKBY, pr.EDIT_INITIALS, pr.LAST_EDIT_COMMENT, pr.MR_FID, pr.dst_start_date from ds_procedure pr
inner join ds_fish fi on fi.f_id = pr.f_id
inner join ds_moriver mo on mo.mr_id = fi.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :1 != 'ZZ' THEN si.fieldoffice ELSE :2 END) = :3`

var procedureDataEntriesCountBySql = `SELECT count(*) from ds_procedure pr
inner join ds_fish fi on fi.f_id = pr.f_id
inner join ds_moriver mo on mo.mr_id = fi.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :1 != 'ZZ' THEN si.fieldoffice ELSE :2 END) = :3`

var procedureDataEntriesByIdSql = `select pr.ID, pr.F_ID, pr.F_FID, si.site_id, pr.PURPOSE_CODE, pr.PROCEDURE_DATE, pr.PROCEDURE_START_TIME, pr.PROCEDURE_END_TIME, pr.PROCEDURE_BY, pr.ANTIBIOTIC_INJECTION_IND, pr.PHOTO_DORSAL_IND, pr.PHOTO_VENTRAL_IND, 
pr.PHOTO_LEFT_IND, pr.OLD_RADIO_TAG_NUM, pr.OLD_FREQUENCY_ID, pr.DST_SERIAL_NUM, pr.DST_START_TIME, pr.DST_REIMPLANT_IND, pr.NEW_RADIO_TAG_NUM, pr.NEW_FREQUENCY_ID, pr.SEX_CODE, pr.COMMENTS, pr.FISH_HEALTH_COMMENTS, pr.SPAWN_CODE, pr.EVAL_LOCATION_CODE, 
pr.BLOOD_SAMPLE_IND, pr.EGG_SAMPLE_IND, pr.VISUAL_REPRO_STATUS_CODE, pr.ULTRASOUND_REPRO_STATUS_CODE, pr.ULTRASOUND_GONAD_LENGTH, pr.GONAD_CONDITION, pr.EXPECTED_SPAWN_YEAR, pr.LAST_UPDATED, pr.UPLOAD_SESSION_ID, pr.UPLOADED_BY, pr.UPLOAD_FILENAME, 
pr.CHECKBY, pr.EDIT_INITIALS, pr.LAST_EDIT_COMMENT, pr.MR_FID, pr.dst_start_date from ds_procedure pr
inner join ds_fish fi on fi.f_id = pr.f_id
inner join ds_moriver mo on mo.mr_id = fi.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4
and pr.id = :1`

var procedureDataEntriesCountByIdSql = `SELECT count(*) from ds_procedure pr
inner join ds_fish fi on fi.f_id = pr.f_id
inner join ds_moriver mo on mo.mr_id = fi.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4
and pr.id = :1`

var procedureDataEntriesByFidSql = `select pr.ID, pr.F_ID, pr.F_FID, si.site_id, pr.PURPOSE_CODE, pr.PROCEDURE_DATE, pr.PROCEDURE_START_TIME, pr.PROCEDURE_END_TIME, pr.PROCEDURE_BY, pr.ANTIBIOTIC_INJECTION_IND, pr.PHOTO_DORSAL_IND, pr.PHOTO_VENTRAL_IND, 
pr.PHOTO_LEFT_IND, pr.OLD_RADIO_TAG_NUM, pr.OLD_FREQUENCY_ID, pr.DST_SERIAL_NUM, pr.DST_START_TIME, pr.DST_REIMPLANT_IND, pr.NEW_RADIO_TAG_NUM, pr.NEW_FREQUENCY_ID, pr.SEX_CODE, pr.COMMENTS, pr.FISH_HEALTH_COMMENTS, pr.SPAWN_CODE, pr.EVAL_LOCATION_CODE, 
pr.BLOOD_SAMPLE_IND, pr.EGG_SAMPLE_IND, pr.VISUAL_REPRO_STATUS_CODE, pr.ULTRASOUND_REPRO_STATUS_CODE, pr.ULTRASOUND_GONAD_LENGTH, pr.GONAD_CONDITION, pr.EXPECTED_SPAWN_YEAR, pr.LAST_UPDATED, pr.UPLOAD_SESSION_ID, pr.UPLOADED_BY, pr.UPLOAD_FILENAME, 
pr.CHECKBY, pr.EDIT_INITIALS, pr.LAST_EDIT_COMMENT, pr.MR_FID, pr.dst_start_date from ds_procedure pr
inner join ds_fish fi on fi.f_id = pr.f_id
inner join ds_moriver mo on mo.mr_id = fi.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4 
and pr.f_id = :1`

var procedureDataEntriesCountByFidSql = `SELECT count(*) from ds_procedure pr
inner join ds_fish fi on fi.f_id = pr.f_id
inner join ds_moriver mo on mo.mr_id = fi.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4 
and pr.f_id = :1`

var procedureDataEntriesByFfidSql = `select pr.ID, pr.F_ID, pr.F_FID, si.site_id, pr.PURPOSE_CODE, pr.PROCEDURE_DATE, pr.PROCEDURE_START_TIME, pr.PROCEDURE_END_TIME, pr.PROCEDURE_BY, pr.ANTIBIOTIC_INJECTION_IND, pr.PHOTO_DORSAL_IND, pr.PHOTO_VENTRAL_IND, 
pr.PHOTO_LEFT_IND, pr.OLD_RADIO_TAG_NUM, pr.OLD_FREQUENCY_ID, pr.DST_SERIAL_NUM, pr.DST_START_TIME, pr.DST_REIMPLANT_IND, pr.NEW_RADIO_TAG_NUM, pr.NEW_FREQUENCY_ID, pr.SEX_CODE, pr.COMMENTS, pr.FISH_HEALTH_COMMENTS, pr.SPAWN_CODE, pr.EVAL_LOCATION_CODE, 
pr.BLOOD_SAMPLE_IND, pr.EGG_SAMPLE_IND, pr.VISUAL_REPRO_STATUS_CODE, pr.ULTRASOUND_REPRO_STATUS_CODE, pr.ULTRASOUND_GONAD_LENGTH, pr.GONAD_CONDITION, pr.EXPECTED_SPAWN_YEAR, pr.LAST_UPDATED, pr.UPLOAD_SESSION_ID, pr.UPLOADED_BY, pr.UPLOAD_FILENAME, 
pr.CHECKBY, pr.EDIT_INITIALS, pr.LAST_EDIT_COMMENT, pr.MR_FID, pr.dst_start_date from ds_procedure pr
inner join ds_fish fi on fi.f_id = pr.f_id
inner join ds_moriver mo on mo.mr_id = fi.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4 
and pr.f_fid = :1`

var procedureDataEntriesCountByFfidSql = `SELECT count(*) from ds_procedure pr
inner join ds_fish fi on fi.f_id = pr.f_id
inner join ds_moriver mo on mo.mr_id = fi.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4 
and pr.f_fid = :1`

var procedureDataEntriesByMrIdSql = `select pr.ID, pr.F_ID, pr.F_FID, si.site_id, pr.PURPOSE_CODE, pr.PROCEDURE_DATE, pr.PROCEDURE_START_TIME, pr.PROCEDURE_END_TIME, pr.PROCEDURE_BY, COALESCE(pr.ANTIBIOTIC_INJECTION_IND,0) as ANTIBIOTIC_INJECTION_IND, COALESCE(pr.PHOTO_DORSAL_IND,0) as PHOTO_DORSAL_IND, COALESCE(pr.PHOTO_VENTRAL_IND,0) as PHOTO_VENTRAL_IND, 
COALESCE(pr.PHOTO_LEFT_IND,0) as PHOTO_LEFT_IND, COALESCE(pr.OLD_RADIO_TAG_NUM,0) as OLD_RADIO_TAG_NUM, COALESCE(pr.OLD_FREQUENCY_ID,0) as OLD_FREQUENCY_ID, COALESCE(pr.DST_SERIAL_NUM,0) as DST_SERIAL_NUM, pr.DST_START_TIME, COALESCE(pr.DST_REIMPLANT_IND,0) as DST_REIMPLANT_IND, COALESCE(pr.NEW_RADIO_TAG_NUM,0) as NEW_RADIO_TAG_NUM,
COALESCE(pr.NEW_FREQUENCY_ID,0) as NEW_FREQUENCY_ID, pr.SEX_CODE, pr.COMMENTS, pr.FISH_HEALTH_COMMENTS, pr.SPAWN_CODE, pr.EVAL_LOCATION_CODE, COALESCE(pr.BLOOD_SAMPLE_IND,0) as BLOOD_SAMPLE_IND, COALESCE(pr.EGG_SAMPLE_IND,0) as EGG_SAMPLE_IND, pr.VISUAL_REPRO_STATUS_CODE, pr.ULTRASOUND_REPRO_STATUS_CODE, COALESCE(pr.ULTRASOUND_GONAD_LENGTH,0) as ULTRASOUND_GONAD_LENGTH, pr.GONAD_CONDITION, 
COALESCE(pr.EXPECTED_SPAWN_YEAR,0) as EXPECTED_SPAWN_YEAR, pr.LAST_UPDATED, pr.UPLOAD_SESSION_ID, pr.UPLOADED_BY, pr.UPLOAD_FILENAME, pr.CHECKBY, pr.EDIT_INITIALS, pr.LAST_EDIT_COMMENT, pr.MR_FID
from ds_procedure pr
inner join ds_fish fi on fi.f_id = pr.f_id
inner join ds_moriver mo on mo.mr_id = fi.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4 
and fi.mr_id = :1`

var procedureDataEntriesCountByMrIdSql = `SELECT count(*) from ds_procedure pr
inner join ds_fish fi on fi.f_id = pr.f_id
inner join ds_moriver mo on mo.mr_id = fi.mr_id
inner join ds_sites si on si.site_id = mo.site_id
where (CASE when :2 != 'ZZ' THEN si.fieldoffice ELSE :3 END) = :4 
and fi.mr_id = :1`

func (s *PallidSturgeonStore) GetProcedureDataEntries(tableId string, fId string, mrId string, officeCode string, queryParams models.SearchParams) (models.ProcedureDataEntryWithCount, error) {
	procedureDataEntryWithCount := models.ProcedureDataEntryWithCount{}
	query := ""
	queryWithCount := ""
	id := ""

	if tableId != "" {
		query = procedureDataEntriesByIdSql
		queryWithCount = procedureDataEntriesCountByIdSql
		id = tableId
	}

	if fId != "" {
		query = procedureDataEntriesByFidSql
		queryWithCount = procedureDataEntriesCountByFidSql
		id = fId
	}

	if mrId != "" {
		query = procedureDataEntriesByMrIdSql
		queryWithCount = procedureDataEntriesCountByMrIdSql
		id = mrId
	}

	if tableId == "" && fId == "" && mrId == "" {
		query = procedureDataEntriesSql
		queryWithCount = procedureDataEntriesCountBySql
	}

	countQuery, err := s.db.Prepare(queryWithCount)
	if err != nil {
		return procedureDataEntryWithCount, err
	}

	var countrows *sql.Rows
	if id == "" {
		countrows, err = countQuery.Query(officeCode, officeCode, officeCode)
		if err != nil {
			return procedureDataEntryWithCount, err
		}
	} else {
		countrows, err = countQuery.Query(officeCode, officeCode, officeCode, id)
		if err != nil {
			return procedureDataEntryWithCount, err
		}
	}
	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&procedureDataEntryWithCount.TotalCount)
		if err != nil {
			return procedureDataEntryWithCount, err
		}
	}

	procedureEntries := []models.UploadProcedure{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "id"
	}
	procedureDataEntriesSqlWithSearch := query + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(procedureDataEntriesSqlWithSearch)
	if err != nil {
		return procedureDataEntryWithCount, err
	}

	var rows *sql.Rows
	if id == "" {
		rows, err = dbQuery.Query(officeCode, officeCode, officeCode)
		if err != nil {
			return procedureDataEntryWithCount, err
		}
	} else {
		rows, err = dbQuery.Query(officeCode, officeCode, officeCode, id)
		if err != nil {
			return procedureDataEntryWithCount, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		procedureDataEntry := models.UploadProcedure{}
		err = rows.Scan(&procedureDataEntry.Id, &procedureDataEntry.Fid, &procedureDataEntry.FFid, &procedureDataEntry.SiteID, &procedureDataEntry.PurposeCode, &procedureDataEntry.ProcedureDate, &procedureDataEntry.ProcedureStartTime, &procedureDataEntry.ProcedureEndTime, &procedureDataEntry.ProcedureBy, &procedureDataEntry.AntibioticInjectionInd,
			&procedureDataEntry.PhotoDorsalInd, &procedureDataEntry.PhotoVentralInd, &procedureDataEntry.PhotoLeftInd, &procedureDataEntry.OldRadioTagNum, &procedureDataEntry.OldFrequencyId, &procedureDataEntry.DstSerialNum, &procedureDataEntry.DstStartTime, &procedureDataEntry.DstReimplantInd,
			&procedureDataEntry.NewRadioTagNum, &procedureDataEntry.NewFrequencyId, &procedureDataEntry.SexCode, &procedureDataEntry.Comments, &procedureDataEntry.FishHealthComments, &procedureDataEntry.SpawnStatus, &procedureDataEntry.EvalLocationCode, &procedureDataEntry.BloodSampleInd, &procedureDataEntry.EggSampleInd,
			&procedureDataEntry.VisualReproStatusCode, &procedureDataEntry.UltrasoundReproStatusCode, &procedureDataEntry.UltrasoundGonadLength, &procedureDataEntry.GonadCondition, &procedureDataEntry.ExpectedSpawnYear, &procedureDataEntry.LastUpdated, &procedureDataEntry.UploadSessionId, &procedureDataEntry.UploadedBy,
			&procedureDataEntry.UploadFilename, &procedureDataEntry.Checkby, &procedureDataEntry.EditInitials, &procedureDataEntry.LastEditComment, &procedureDataEntry.MrFid, &procedureDataEntry.DstStartDate)
		if err != nil {
			return procedureDataEntryWithCount, err
		}
		procedureEntries = append(procedureEntries, procedureDataEntry)
	}

	procedureDataEntryWithCount.Items = procedureEntries

	return procedureDataEntryWithCount, err
}

var insertProcedureDataSql = `insert into ds_procedure (ID, F_ID, F_FID, PURPOSE_CODE, PROCEDURE_DATE, PROCEDURE_START_TIME, PROCEDURE_END_TIME, PROCEDURE_BY, ANTIBIOTIC_INJECTION_IND, PHOTO_DORSAL_IND, PHOTO_VENTRAL_IND, PHOTO_LEFT_IND, OLD_RADIO_TAG_NUM, OLD_FREQUENCY_ID, DST_SERIAL_NUM, DST_START_DATE, DST_START_TIME, 
	DST_REIMPLANT_IND, NEW_RADIO_TAG_NUM, NEW_FREQUENCY_ID, SEX_CODE, COMMENTS, FISH_HEALTH_COMMENTS, SPAWN_CODE, EVAL_LOCATION_CODE, BLOOD_SAMPLE_IND, EGG_SAMPLE_IND, VISUAL_REPRO_STATUS_CODE, ULTRASOUND_REPRO_STATUS_CODE, ULTRASOUND_GONAD_LENGTH, GONAD_CONDITION, EXPECTED_SPAWN_YEAR, LAST_UPDATED, UPLOAD_SESSION_ID, UPLOADED_BY, 
	UPLOAD_FILENAME, CHECKBY, EDIT_INITIALS, LAST_EDIT_COMMENT, MR_FID) 
	values (procedure_seq.nextval, :1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23,:24,:25,:26,:27,:28,:29,:30,:31,:32,:33,:34,:35,:36,:37,:38,:39)`

func (s *PallidSturgeonStore) SaveProcedureDataEntry(procedureDataEntry models.UploadProcedure) error {
	_, err := s.db.Exec(insertProcedureDataSql, procedureDataEntry.Fid, procedureDataEntry.FFid, procedureDataEntry.PurposeCode, procedureDataEntry.ProcedureDate, procedureDataEntry.ProcedureStartTime, procedureDataEntry.ProcedureEndTime, procedureDataEntry.ProcedureBy, procedureDataEntry.AntibioticInjectionInd,
		procedureDataEntry.PhotoDorsalInd, procedureDataEntry.PhotoVentralInd, procedureDataEntry.PhotoLeftInd, procedureDataEntry.OldRadioTagNum, procedureDataEntry.OldFrequencyId, procedureDataEntry.DstSerialNum, procedureDataEntry.DstStartDate, procedureDataEntry.DstStartTime, procedureDataEntry.DstReimplantInd,
		procedureDataEntry.NewRadioTagNum, procedureDataEntry.NewFrequencyId, procedureDataEntry.SexCode, procedureDataEntry.Comments, procedureDataEntry.FishHealthComments, procedureDataEntry.SpawnStatus, procedureDataEntry.EvalLocationCode, procedureDataEntry.BloodSampleInd, procedureDataEntry.EggSampleInd,
		procedureDataEntry.VisualReproStatusCode, procedureDataEntry.UltrasoundReproStatusCode, procedureDataEntry.UltrasoundGonadLength, procedureDataEntry.GonadCondition, procedureDataEntry.ExpectedSpawnYear, procedureDataEntry.LastUpdated, procedureDataEntry.UploadSessionId, procedureDataEntry.UploadedBy,
		procedureDataEntry.UploadFilename, procedureDataEntry.Checkby, procedureDataEntry.EditInitials, procedureDataEntry.LastEditComment, procedureDataEntry.MrFid)
	return err
}

var updateProcedureDataSql = `update ds_procedure set 
f_id = :2,
f_fid =:3,
PURPOSE_CODE = :4, 
PROCEDURE_DATE = :5, 
PROCEDURE_START_TIME = :6, 
PROCEDURE_END_TIME = :7, 
PROCEDURE_BY = :8,
ANTIBIOTIC_INJECTION_IND = :9, 
PHOTO_DORSAL_IND = :10, 
PHOTO_VENTRAL_IND = :11, 
PHOTO_LEFT_IND = :12, 
OLD_RADIO_TAG_NUM = :13,
OLD_FREQUENCY_ID = :14, 
DST_SERIAL_NUM = :15, 
DST_START_DATE = :16, 
DST_START_TIME = :17, 
DST_REIMPLANT_IND = :18, 
NEW_RADIO_TAG_NUM = :19, 
NEW_FREQUENCY_ID = :20, 
SEX_CODE = :21, 
COMMENTS = :22, 
FISH_HEALTH_COMMENTS = :23, 
SPAWN_CODE = :24, 
EVAL_LOCATION_CODE = :25, 
BLOOD_SAMPLE_IND = :26,
EGG_SAMPLE_IND = :27, 
VISUAL_REPRO_STATUS_CODE = :28, 
ULTRASOUND_REPRO_STATUS_CODE = :29, 
ULTRASOUND_GONAD_LENGTH = :30, 
GONAD_CONDITION = :31, 
EXPECTED_SPAWN_YEAR = :32, 
LAST_UPDATED = :33, 
UPLOAD_SESSION_ID = :34, 
UPLOADED_BY = :35, 
UPLOAD_FILENAME = :36, 
CHECKBY = :37,
EDIT_INITIALS = :38, 
LAST_EDIT_COMMENT = :39, 
mr_fid = :40
where id = :1`

func (s *PallidSturgeonStore) UpdateProcedureDataEntry(procedureDataEntry models.UploadProcedure) error {
	_, err := s.db.Exec(updateProcedureDataSql, procedureDataEntry.Fid, procedureDataEntry.FFid, procedureDataEntry.PurposeCode, procedureDataEntry.ProcedureDate, procedureDataEntry.ProcedureStartTime, procedureDataEntry.ProcedureEndTime, procedureDataEntry.ProcedureBy, procedureDataEntry.AntibioticInjectionInd,
		procedureDataEntry.PhotoDorsalInd, procedureDataEntry.PhotoVentralInd, procedureDataEntry.PhotoLeftInd, procedureDataEntry.OldRadioTagNum, procedureDataEntry.OldFrequencyId, procedureDataEntry.DstSerialNum, procedureDataEntry.DstStartDate, procedureDataEntry.DstStartTime, procedureDataEntry.DstReimplantInd,
		procedureDataEntry.NewRadioTagNum, procedureDataEntry.NewFrequencyId, procedureDataEntry.SexCode, procedureDataEntry.Comments, procedureDataEntry.FishHealthComments, procedureDataEntry.SpawnStatus, procedureDataEntry.EvalLocationCode, procedureDataEntry.BloodSampleInd, procedureDataEntry.EggSampleInd,
		procedureDataEntry.VisualReproStatusCode, procedureDataEntry.UltrasoundReproStatusCode, procedureDataEntry.UltrasoundGonadLength, procedureDataEntry.GonadCondition, procedureDataEntry.ExpectedSpawnYear, procedureDataEntry.LastUpdated, procedureDataEntry.UploadSessionId, procedureDataEntry.UploadedBy,
		procedureDataEntry.UploadFilename, procedureDataEntry.Checkby, procedureDataEntry.EditInitials, procedureDataEntry.LastEditComment, procedureDataEntry.MrFid, procedureDataEntry.Id)
	return err
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

var fishDataSummarySql = `SELECT mr_id, f_id, year, FIELD_OFFICE_CODE, PROJECT_CODE, SEGMENT_CODE, SEASON_CODE, COALESCE(BEND_NUMBER, 0) as BEND_NUMBER, BEND_R_OR_N, COALESCE(bend_river_mile, 0.0) as bend_river_mile, panelhook, SPECIES_CODE, HATCHERY_ORIGIN_CODE, checkby FROM table (pallid_data_api.fish_datasummary_fnc(:1, :2, :3, :4, :5, :6, :7, to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

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

var missouriDataSummarySql = `SELECT mr_id, year, FIELD_OFFICE_CODE, PROJECT_CODE, SEGMENT_CODE, SEASON_CODE, COALESCE(BEND_NUMBER,0) as bend_number, BEND_R_OR_N, COALESCE(bend_river_mile, 0.0) as bend_river_mile, COALESCE(subsample,0) as subsample, COALESCE(subsample_pass,0) as subsample_pass, set_Date, conductivity, checkby, COALESCE(approved, 0) as approved FROM table (pallid_data_api.missouri_datasummary_fnc(:1, :2, :3, :4, :5, :6, :7, to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

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
			&summary.Subsample, &summary.Pass, &summary.SetDate, &summary.Conductivity, &summary.CheckedBy, &summary.Approved)
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

var geneticDataSummarySql = `SELECT year,FIELD_OFFICE_CODE,PROJECT_CODE,genetics_vial_number,pit_tag,river,COALESCE(river_mile, 0.0) as river_mile,
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

var searchDataSummaryFullDataSql = `SELECT * FROM table (pallid_data_api.search_datasummary_fnc(:1,:2,:3,:4,:5,:6,:7,to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

func (s *PallidSturgeonStore) GetFullSearchDataSummary(year string, officeCode string, project string, approved string, season string, segment string, month string, fromDate string, toDate string) (string, error) {
	dbQuery, err := s.db.Prepare(searchDataSummaryFullDataSql)
	if err != nil {
		return "Cannot create file", err
	}

	rows, err := dbQuery.Query(year, officeCode, project, approved, season, segment, month, fromDate, toDate)
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

var searchDataSummarySql = `SELECT year,fieldoffice,project_id,segment_id,season,se_id,search_date,recorder,search_type_code,start_time,start_latitude,start_longitude,stop_time,stop_latitude,stop_longitude,temp,conductivity
FROM table (pallid_data_api.search_datasummary_fnc(:1,:2,:3,:4,:5,:6,:7,to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

var searchDataSummaryCountSql = `SELECT count(*) FROM table (pallid_data_api.search_datasummary_fnc(:1,:2,:3,:4,:5,:6,:7,to_date(:8,'MM/DD/YYYY'), to_date(:9,'MM/DD/YYYY')))`

func (s *PallidSturgeonStore) GetSearchDataSummary(year string, officeCode string, project string, approved string, season string, segment string, month string, fromDate string, toDate string, queryParams models.SearchParams) (models.SearchSummaryWithCount, error) {
	searchSummariesWithCount := models.SearchSummaryWithCount{}
	countQuery, err := s.db.Prepare(searchDataSummaryCountSql)
	if err != nil {
		return searchSummariesWithCount, err
	}

	countrows, err := countQuery.Query(year, officeCode, project, approved, season, segment, month, fromDate, toDate)
	if err != nil {
		return searchSummariesWithCount, err
	}
	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&searchSummariesWithCount.TotalCount)
		if err != nil {
			return searchSummariesWithCount, err
		}
	}

	searchSummaries := []models.SearchSummary{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "se_id"
	}
	searchDataSummarySqlWithSearch := searchDataSummarySql + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(searchDataSummarySqlWithSearch)
	if err != nil {
		return searchSummariesWithCount, err
	}

	rows, err := dbQuery.Query(year, officeCode, project, approved, season, segment, month, fromDate, toDate)
	if err != nil {
		return searchSummariesWithCount, err
	}
	defer rows.Close()

	for rows.Next() {
		summary := models.SearchSummary{}
		err = rows.Scan(&summary.Year, &summary.FieldOffice, &summary.Project, &summary.Segment, &summary.Season, &summary.SeID, &summary.SearchDate, &summary.Recorder, &summary.SearchTypeCode, &summary.StartTime,
			&summary.StartLatitude, &summary.StartLongitude, &summary.StopTime, &summary.StopLatitude, &summary.StopLongitude, &summary.Temp, &summary.Conductivity)
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

var missouriDatasheetsBySiteId = `select site_id, mr_id, mr_fid, subsample, subsamplepass, subsamplen, recorder, conductivity, bkg_color, fish_count, supp_count, supp_bkg_color from table (pallid_data_entry_api.data_entry_missouri_fnc(:1,:2,:3,:4,:5,:6))`

var missouriDatasheetsCountBySiteId = `select count(*) from table (pallid_data_entry_api.data_entry_missouri_fnc(:1,:2,:3,:4,:5,:6))`

func (s *PallidSturgeonStore) GetMissouriDatasheetById(siteId string, officeCode string, project string, segment string, season string, bend string, queryParams models.SearchParams) (models.MoriverDataEntryWithCount, error) {
	missouriDatasheetsWithCount := models.MoriverDataEntryWithCount{}
	countQuery, err := s.db.Prepare(missouriDatasheetsCountBySiteId)
	if err != nil {
		return missouriDatasheetsWithCount, err
	}

	countrows, err := countQuery.Query(siteId, officeCode, project, segment, season, bend)
	if err != nil {
		return missouriDatasheetsWithCount, err
	}
	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&missouriDatasheetsWithCount.TotalCount)
		if err != nil {
			return missouriDatasheetsWithCount, err
		}
	}

	missouriDatasheets := []models.UploadMoriver{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "site_id"
	}
	missouriDataByIdSqlWithSearch := missouriDatasheetsBySiteId + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(missouriDataByIdSqlWithSearch)
	if err != nil {
		return missouriDatasheetsWithCount, err
	}

	rows, err := dbQuery.Query(siteId, officeCode, project, segment, season, bend)
	if err != nil {
		return missouriDatasheetsWithCount, err
	}
	defer rows.Close()

	for rows.Next() {
		datasheets := models.UploadMoriver{}
		err = rows.Scan(&datasheets.SiteID, &datasheets.MrID, &datasheets.MrFid, &datasheets.Subsample, &datasheets.Subsamplepass, &datasheets.Subsamplen, &datasheets.Recorder, &datasheets.Conductivity, &datasheets.BkgColor,
			&datasheets.FishCount, &datasheets.SuppCount, &datasheets.SuppBkgColor)
		if err != nil {
			return missouriDatasheetsWithCount, err
		}
		missouriDatasheets = append(missouriDatasheets, datasheets)
	}

	missouriDatasheetsWithCount.Items = missouriDatasheets

	return missouriDatasheetsWithCount, err
}

var searchDatasheetsBySiteId = `select si.site_id, se.se_id, se.recorder, se.search_type_code, se.start_time, se.start_latitude, se.start_longitude, se.stop_time, se.stop_latitude, se.stop_longitude, 
se.temp, se.conductivity from ds_sites si inner join ds_search se on si.site_id = se.site_id where si.site_id = :1`

var searchDatasheetsCountBySiteId = `select count(*) from ds_sites si inner join ds_search se on se.site_id = si.site_id where si.site_id = :1`

func (s *PallidSturgeonStore) GetSearchDatasheetById(siteId string, queryParams models.SearchParams) (models.UploadSearchData, error) {
	searchDatasheetsWithCount := models.UploadSearchData{}
	countQuery, err := s.db.Prepare(searchDatasheetsCountBySiteId)
	if err != nil {
		return searchDatasheetsWithCount, err
	}

	countrows, err := countQuery.Query(siteId)
	if err != nil {
		return searchDatasheetsWithCount, err
	}
	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&searchDatasheetsWithCount.TotalCount)
		if err != nil {
			return searchDatasheetsWithCount, err
		}
	}

	searchDatasheets := []models.UploadSearch{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "site_id"
	}
	sqlQueryWithSearch := searchDatasheetsBySiteId + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(sqlQueryWithSearch)
	if err != nil {
		return searchDatasheetsWithCount, err
	}

	rows, err := dbQuery.Query(siteId)
	if err != nil {
		return searchDatasheetsWithCount, err
	}
	defer rows.Close()

	for rows.Next() {
		datasheets := models.UploadSearch{}
		err = rows.Scan(&datasheets.SiteId, &datasheets.SeId, &datasheets.Recorder, &datasheets.SearchTypeCode, &datasheets.StartTime, &datasheets.StartLatitude, &datasheets.StartLongitude, &datasheets.StopTime,
			&datasheets.StopLatitude, &datasheets.StopLongitude, &datasheets.Temp, &datasheets.Conductivity)
		if err != nil {
			return searchDatasheetsWithCount, err
		}
		searchDatasheets = append(searchDatasheets, datasheets)
	}

	searchDatasheetsWithCount.Items = searchDatasheets

	return searchDatasheetsWithCount, err
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
	segment_id, segment, season_id, season, bend, bendrn, bend_river_mile, comments, last_updated, upload_session_id, uploaded_by, upload_filename)
	values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19)`

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
		uploadSite.LastUpdated,
		uploadSite.UploadSessionId,
		uploadSite.UploadedBy,
		uploadSite.UploadFilename,
	)

	return err
}

var insertFishUploadSql = `insert into upload_fish (site_id, f_fid, mr_fid, panelhook, bait, species, length, weight,
	fishcount, fin_curl, otolith, rayspine, scale, ftprefix, ftnum, ftmr,
	comments, last_updated, upload_session_id, uploaded_by, upload_filename)
	 values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21)`

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
		uploadFish.LastUpdated,
		uploadFish.UploadSessionId,
		uploadFish.UploadedBy,
		uploadFish.UploadFilename,
	)

	return err
}

var insertSearchUploadSql = `insert into upload_search(se_fid, ds_id, site_id, site_fid, search_date, recorder, search_type_code, search_day, start_time,  
		start_latitude, start_longitude, stop_time, stop_latitude, stop_longitude, temp, conductivity, last_updated, 
		upload_session_id, uploaded_by, upload_filename)
	values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20)`

func (s *PallidSturgeonStore) SaveSearchUpload(uploadSearch models.UploadSearch) error {
	_, err := s.db.Exec(insertSearchUploadSql,
		uploadSearch.SeFid,
		uploadSearch.DsId,
		uploadSearch.SiteId,
		uploadSearch.SiteFid,
		uploadSearch.SearchDate,
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
	comments,
	last_updated, upload_session_id, uploaded_by, upload_filename) values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,
		:21,:22,:23,:24,:25,:26,:27,:28,:29,:30,:31,:32,:33,:34,:35,:36,:37,:38,:39,:40,:41,:42,:43,:44,:45)`

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
	expected_spawn_year, ultrasound_gonad_length, gonad_condition, last_updated, upload_session_id, uploaded_by, upload_filename)                                                        
values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23,:24,:25,:26,:27,:28,:29,:30,:31,:32,:33,:34,:35)`

func (s *PallidSturgeonStore) SaveProcedureUpload(uploadProcedure models.UploadProcedure) error {
	_, err := s.db.Exec(insertProcedureUploadSql,
		uploadProcedure.FFid,
		uploadProcedure.MrFid,
		uploadProcedure.PurposeCode,
		uploadProcedure.ProcedureDate,
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
		uploadProcedure.DstStartDate,
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
	comments, last_updated, upload_session_id,
	uploaded_by, upload_filename, complete,
	no_turbidity, no_velocity)

	 values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23,:24,
		:25,:26,:27,:28,:29,:30,:31,:32,:33,:34,:35,:36,:37,:38,:39,:40,:41,:42,:43,:44,:45,:46,:47,
		:48,:49,:50,:51,:52,:53,:54,:55,:56,:57,:58,:59,:60,:61,:62,:63,:64,:65,:66,:67,:68,:69,:70,
		:71)`

func (s *PallidSturgeonStore) SaveMoriverUpload(UploadMoriver models.UploadMoriver) error {
	_, err := s.db.Exec(insertMoriverUploadSql,
		UploadMoriver.SiteID, UploadMoriver.SiteFid, UploadMoriver.MrFid, UploadMoriver.Season, UploadMoriver.SetDate,
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
		UploadMoriver.Comments, UploadMoriver.LastUpdated, UploadMoriver.UploadSessionId,
		UploadMoriver.UploadedBy, UploadMoriver.UploadFilename, UploadMoriver.Complete,
		UploadMoriver.NoTurbidity, UploadMoriver.NoVelocity,
	)

	return err
}

var insertTelemetryUploadSql = `insert into upload_telemetry_fish(t_fid, se_fid, bend, radio_tag_num, frequency_id_code, capture_time, capture_latitude, capture_longitude,
	position_confidence, macro_id, meso_id, depth, temp, conductivity, turbidity, silt, sand, gravel, comments, last_updated, upload_session_id, uploaded_by, upload_filename)
	values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23)`

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

var errorCountSql = `select el.year, count(el.el_id)
from site_error_log_v el
where NVL(error_fixed,0) = 0
and (case
when el.worksheet_type_id = 2 
then (select FIELDOFFICE
from ds_sites
where site_id = (select site_id
from ds_moriver
where mr_id = el.worksheet_id))
when el.worksheet_type_id = 1 
then NULL
when el.worksheet_type_id in(3,4) 
then (select FIELDOFFICE
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

var getOfficeErrorLogSql = `select el.site_id, el.year, el.el_id, el.error_entry_date, el.worksheet_type_id, el.field_id, el.error_description, COALESCE(el.error_fixed, 0) as error_fixed, el.worksheet_id, el.form_id
from site_error_log_v el
where (case
when el.worksheet_type_id = 2 
then (select FIELDOFFICE
from ds_sites
where site_id = (select site_id
from ds_moriver
where mr_id = el.worksheet_id))
when el.worksheet_type_id = 1 
then NULL
when el.worksheet_type_id in(3,4) 
then (select FIELDOFFICE
from ds_sites
where site_id = (select mr2.site_id
from ds_sites s2, ds_moriver mr2, ds_fish f2
where s2.site_id = mr2.site_id
and mr2.mr_id = F2.MR_ID
and f2.f_id = el.worksheet_id))
end) = :1`

func (s *PallidSturgeonStore) GetOfficeErrorLogs(fieldOfficeCode string) ([]models.OfficeErrorLog, error) {
	officeErrorLogs := []models.OfficeErrorLog{}

	rows, err := s.db.Query(getOfficeErrorLogSql, fieldOfficeCode)
	if err != nil {
		return officeErrorLogs, err
	}
	defer rows.Close()

	for rows.Next() {
		officeErrorLog := models.OfficeErrorLog{}
		err = rows.Scan(&officeErrorLog.SiteID, &officeErrorLog.Year, &officeErrorLog.ElID, &officeErrorLog.ErrorEntryDate, &officeErrorLog.WorksheetTypeID, &officeErrorLog.FieldID,
			&officeErrorLog.ErrorDescription, &officeErrorLog.ErrorStatus, &officeErrorLog.WorksheetID, &officeErrorLog.FormID)
		if err != nil {
			return officeErrorLogs, err
		}
		officeErrorLogs = append(officeErrorLogs, officeErrorLog)
	}

	return officeErrorLogs, err
}

var usgNoVialNumberSql = `select fo.field_office_description||' : '||p.project_description as fp,
f.species, 
f.f_id, mr.mr_id, MR.SITE_ID as mrsite_id,   DS.SITE_ID as s_site_id,
f.f_fid, Sup.GENETICS_VIAL_NUMBER
from ds_fish f, ds_supplemental sup, ds_moriver mr, ds_sites ds, project_lk p, segment_lk s, field_office_lk fo
where F.F_ID = Sup.F_ID (+)
and MR.MR_ID = F.MR_ID (+)
and mr.site_id = ds.site_id (+)
and DS.PROJECT_ID = P.PROJECT_CODE (+)
and DS.FIELDOFFICE = fo.FIELD_OFFICE_CODE (+)
and ds.SEGMENT_ID = s.segment_code (+)
and (f.species = 'USG' or f.species = 'PDSG')
and Sup.GENETICS_VIAL_NUMBER IS NULL
and (CASE when :1 != 'ZZ' THEN ds.FIELDOFFICE ELSE :2 END) = :3
and ds.PROJECT_ID = :4
order by ds.FIELDOFFICE, ds.PROJECT_ID, ds.SEGMENT_ID, ds.BEND`

func (s *PallidSturgeonStore) GetUsgNoVialNumbers(fieldOfficeCode string, projectCode string) ([]models.UsgNoVialNumber, error) {
	usgNoVialNumbers := []models.UsgNoVialNumber{}

	rows, err := s.db.Query(usgNoVialNumberSql, fieldOfficeCode, fieldOfficeCode, fieldOfficeCode, projectCode)
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

var unapprovedDataSheetsSql = `select asv.ch,
f.field_office_description||' : '||p.project_description as fp, 
s.segment_description,
COALESCE(m.bend, 0) as bend,
m.MR_ID, 
m.SUBSAMPLE,
m.RECORDER,
m.CHECKBY,
COALESCE(m.NETRIVERMILE, 0) as NETRIVERMILE,
m.site_id,
ds.project_id, ds.segment_id, ds.season, ds.fieldoffice,
ds.sample_unit_type,
m.gear
from DS_MORIVER m, project_lk p, segment_lk s, field_office_lk f, approval_status_v asv, ds_sites ds
where m.site_id = ds.site_id (+)
and ds.segment_id = s.segment_code (+)
and DS.PROJECT_ID = P.PROJECT_CODE (+)
and DS.FIELDOFFICE = F.FIELD_OFFICE_CODE
and m.mr_id = asv.mr_id (+)
and asv.ch = 'Unapproved'
and asv.cb = 'YES'
and ds.project_id = :1
and (CASE when :2 != 'ZZ' THEN ds.FIELDOFFICE ELSE :3 END) = :4`

var unapprovedDataSheetsCountSql = `select count(*)
from DS_MORIVER m, project_lk p, segment_lk s, field_office_lk f, approval_status_v asv, ds_sites ds
where m.site_id = ds.site_id (+)
and ds.segment_id = s.segment_code (+)
and DS.PROJECT_ID = P.PROJECT_CODE (+)
and DS.FIELDOFFICE = F.FIELD_OFFICE_CODE
and m.mr_id = asv.mr_id (+)
and asv.ch = 'Unapproved'
and asv.cb = 'YES'
and ds.project_id = :1
and (CASE when :2 != 'ZZ' THEN ds.FIELDOFFICE ELSE :3 END) = :4`

func (s *PallidSturgeonStore) GetUnapprovedDataSheets(projectCode string, officeCode string, queryParams models.SearchParams) (models.UnapprovedDataWithCount, error) {
	unapprovedDataSheetsWithCount := models.UnapprovedDataWithCount{}
	countQuery, err := s.db.Prepare(unapprovedDataSheetsCountSql)
	if err != nil {
		return unapprovedDataSheetsWithCount, err
	}

	countrows, err := countQuery.Query(projectCode, officeCode, officeCode, officeCode)
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

	unapprovedDataSheets := []models.UnapprovedData{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "mr_id"
	}
	selectQueryWithSearch := unapprovedDataSheetsSql + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(selectQueryWithSearch)
	if err != nil {
		return unapprovedDataSheetsWithCount, err
	}

	rows, err := dbQuery.Query(projectCode, officeCode, officeCode, officeCode)
	if err != nil {
		return unapprovedDataSheetsWithCount, err
	}
	defer rows.Close()

	for rows.Next() {
		unapprovedData := models.UnapprovedData{}
		err = rows.Scan(&unapprovedData.Ch, &unapprovedData.Fp, &unapprovedData.SegmentDescription, &unapprovedData.Bend, &unapprovedData.MrId, &unapprovedData.Subsample,
			&unapprovedData.Recorder, &unapprovedData.Checkby, &unapprovedData.NetRiverMile, &unapprovedData.SiteId, &unapprovedData.ProjectId, &unapprovedData.SegmentId, &unapprovedData.Season,
			&unapprovedData.FieldOffice, &unapprovedData.SampleUnitType, &unapprovedData.Gear)
		if err != nil {
			return unapprovedDataSheetsWithCount, err
		}
		unapprovedDataSheets = append(unapprovedDataSheets, unapprovedData)
	}

	unapprovedDataSheetsWithCount.Items = unapprovedDataSheets

	return unapprovedDataSheetsWithCount, err
}

var bafiDataSheetsSql = `SELECT p.project_description||' : '||s.segment_description||' : Bend '||ds.bend as psb,
ds.site_id, DS.FIELDOFFICE, f.f_id, mr.mr_id, mr.mr_fid, F.SPECIES,
MR.RECORDER, MR.SUBSAMPLE, MR.GEAR, F.FISHCOUNT,
ds.year, ds.segment_id, ds.bend, ds.bendrn, 
COALESCE(mr.bendrivermile, 0) as bendrivermile, f.panelhook
from ds_sites ds, ds_moriver mr, ds_fish f, project_lk p, segment_lk s
where DS.SITE_ID = MR.SITE_ID (+)
and MR.MR_ID = F.MR_ID (+)
and DS.PROJECT_ID = P.PROJECT_CODE (+)
and ds.segment_id = s.segment_code (+)
and F.SPECIES = 'BAFI'
and ds.project_id = :1
and (CASE when :2 != 'ZZ' THEN ds.FIELDOFFICE ELSE :3 END) = :4`

var bafiDataSheetCountsSql = `SELECT count(*)
from ds_sites ds, ds_moriver mr, ds_fish f, project_lk p, segment_lk s
where DS.SITE_ID = MR.SITE_ID (+)
and MR.MR_ID = F.MR_ID (+)
and DS.PROJECT_ID = P.PROJECT_CODE (+)
and ds.segment_id = s.segment_code (+)
and F.SPECIES = 'BAFI'
and ds.project_id = :1
and (CASE when :2 != 'ZZ' THEN ds.FIELDOFFICE ELSE :3 END) = :4`

func (s *PallidSturgeonStore) GetBafiDataSheets(fieldOffice string, projectCode string, queryParams models.SearchParams) (models.BafiDataWithCount, error) {
	bafiDataSheetsWithCount := models.BafiDataWithCount{}
	countQuery, err := s.db.Prepare(bafiDataSheetCountsSql)
	if err != nil {
		return bafiDataSheetsWithCount, err
	}

	countrows, err := countQuery.Query(projectCode, fieldOffice, fieldOffice, fieldOffice)
	if err != nil {
		return bafiDataSheetsWithCount, err
	}
	defer countrows.Close()

	for countrows.Next() {
		err = countrows.Scan(&bafiDataSheetsWithCount.TotalCount)
		if err != nil {
			return bafiDataSheetsWithCount, err
		}
	}

	bafiDataSheets := []models.BafiData{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "mr_id"
	}
	selectQueryWithSearch := bafiDataSheetsSql + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(selectQueryWithSearch)
	if err != nil {
		return bafiDataSheetsWithCount, err
	}

	rows, err := dbQuery.Query(projectCode, fieldOffice, fieldOffice, fieldOffice)
	if err != nil {
		return bafiDataSheetsWithCount, err
	}
	defer rows.Close()

	for rows.Next() {
		bafiData := models.BafiData{}
		err = rows.Scan(&bafiData.Psb, &bafiData.SiteId, &bafiData.FieldOffice, &bafiData.FId, &bafiData.MrId, &bafiData.MrFid, &bafiData.Species, &bafiData.Recorder, &bafiData.Subsample, &bafiData.Gear, &bafiData.FishCount,
			&bafiData.Year, &bafiData.SegmentId, &bafiData.Bend, &bafiData.Bendrn, &bafiData.BendRiverMile, &bafiData.PanelHook)
		if err != nil {
			return bafiDataSheetsWithCount, err
		}
		bafiDataSheets = append(bafiDataSheets, bafiData)
	}

	bafiDataSheetsWithCount.Items = bafiDataSheets

	return bafiDataSheetsWithCount, err
}

var uncheckedDataSheetsSql = `select asv.cb,
p.project_description||' : '||s.segment_description||' : Bend '||ds.BEND as psb,
m.MR_ID, 
m.SUBSAMPLE,
m.RECORDER,
m.CHECKBY,
COALESCE(m.NETRIVERMILE, 0) as NETRIVERMILE,
m.site_id,
ds.PROJECT_ID, ds.SEGMENT_ID, ds.SEASON, ds.FIELDOFFICE, m.gear
from DS_MORIVER m, project_lk p, segment_lk s, approval_status_v asv, ds_sites ds
where m.site_id = ds.site_id (+)
and ds.SEGMENT_ID = s.segment_code (+)
and DS.PROJECT_ID = P.project_code (+)
and m.mr_id = asv.mr_id (+)  
and asv.cb = 'Unchecked'
and ds.PROJECT_ID = :1
and M.MR_ID NOT IN (SELECT MR_ID 
FROM DS_FISH
WHERE SPECIES = 'BAFI')
and (CASE when :2 != 'ZZ' THEN ds.FIELDOFFICE ELSE :3 END) = :4`

var uncheckedDataSheetsCountSql = `select count(*)
from DS_MORIVER m, project_lk p, segment_lk s, approval_status_v asv, ds_sites ds
where m.site_id = ds.site_id (+)
and ds.segment_id = s.segment_code (+)
and DS.PROJECT_ID = p.project_code (+)
and m.mr_id = asv.mr_id (+)  
and asv.cb = 'Unchecked'
and ds.PROJECT_ID = :1
and M.MR_ID NOT IN (SELECT MR_ID 
FROM DS_FISH
WHERE SPECIES = 'BAFI')
and (CASE when :2 != 'ZZ' THEN ds.FIELDOFFICE ELSE :3 END) = :4`

func (s *PallidSturgeonStore) GetUncheckedDataSheets(fieldOfficeCode string, projectCode string, queryParams models.SearchParams) (models.UncheckedDataWithCount, error) {
	uncheckedDataSheetsWithCount := models.UncheckedDataWithCount{}
	countQuery, err := s.db.Prepare(uncheckedDataSheetsCountSql)
	if err != nil {
		return uncheckedDataSheetsWithCount, err
	}

	countrows, err := countQuery.Query(projectCode, fieldOfficeCode, fieldOfficeCode, fieldOfficeCode)
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

	uncheckedDataSheets := []models.UncheckedData{}
	offset := queryParams.PageSize * queryParams.Page
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "project_id"
	}
	selectQueryWithSearch := uncheckedDataSheetsSql + fmt.Sprintf(" order by %s OFFSET %s ROWS FETCH NEXT %s ROWS ONLY", queryParams.OrderBy, strconv.Itoa(offset), strconv.Itoa(queryParams.PageSize))
	dbQuery, err := s.db.Prepare(selectQueryWithSearch)
	if err != nil {
		return uncheckedDataSheetsWithCount, err
	}

	rows, err := dbQuery.Query(projectCode, fieldOfficeCode, fieldOfficeCode, fieldOfficeCode)
	if err != nil {
		return uncheckedDataSheetsWithCount, err
	}
	defer rows.Close()

	for rows.Next() {
		uncheckedData := models.UncheckedData{}
		err = rows.Scan(&uncheckedData.Cb, &uncheckedData.Psb, &uncheckedData.MrID, &uncheckedData.Subsample, &uncheckedData.Recorder, &uncheckedData.Checkby, &uncheckedData.Netrivermile, &uncheckedData.SiteID,
			&uncheckedData.ProjectID, &uncheckedData.SegmentID, &uncheckedData.Season, &uncheckedData.FieldOffice, &uncheckedData.Gear)
		if err != nil {
			return uncheckedDataSheetsWithCount, err
		}
		uncheckedDataSheets = append(uncheckedDataSheets, uncheckedData)
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

var getUploadSessionLogsSql = `select debug_text, date_created, p_user, upload_session_id from upload_session_log_t where p_user = :1 and upload_session_id = :2`

func (s *PallidSturgeonStore) GetUploadSessionLogs(user string, uploadSessionId string) ([]models.UploadSessionLog, error) {
	rows, err := s.db.Query(getUploadSessionLogsSql, user, uploadSessionId)

	logs := []models.UploadSessionLog{}
	if err != nil {
		return logs, err
	}
	defer rows.Close()

	for rows.Next() {
		log := models.UploadSessionLog{}
		err = rows.Scan(&log.DebugText, &log.DateCreated, &log.PUser, &log.UploadSessionId)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, err
}

var getSitesExportSql = `select site_id,COALESCE(site_fid, 0) as site_fid,year,fieldoffice,field_office_description,project_id,project_description,segment_id,segment_description,season,season_description,sample_unit_type,bend,bendrn,
COALESCE(bend_river_mile, 0) as bend_river_mile,sample_unit_desc from table (pallid_data_entry_api.data_entry_site_fnc(:1,:2,:3,:4,:5,:6))`

func (s *PallidSturgeonStore) GetSitesExport(year string, officeCode string, project string, segment string, season string, bendrn string) ([]models.ExportSite, error) {
	rows, err := s.db.Query(getSitesExportSql, year, officeCode, project, bendrn, season, segment)

	exportData := []models.ExportSite{}
	if err != nil {
		return exportData, err
	}
	defer rows.Close()

	for rows.Next() {
		export := models.ExportSite{}
		err = rows.Scan(&export.SiteID, &export.SiteFID, &export.SiteYear, &export.FieldOfficeID, &export.FieldOffice, &export.ProjectId, &export.Project, &export.SegmentId, &export.Segment, &export.SeasonId, &export.Season,
			&export.SampleUnitType, &export.Bend, &export.Bendrn, &export.BendRiverMile, &export.SampleUnitDesc)
		if err != nil {
			return nil, err
		}
		exportData = append(exportData, export)
	}

	return exportData, err
}
