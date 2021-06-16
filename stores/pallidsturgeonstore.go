package stores

import (
	"database/sql"

	"di2e.net/cwbi/pallid_sturgeon_api/server/config"
	"di2e.net/cwbi/pallid_sturgeon_api/server/models"
	"github.com/godror/godror"
)

type PallidSturgeonStore struct {
	db     *sql.DB
	config *config.AppConfig
}

var seasonsSql = "select * from season_lk order by id"

func (s *PallidSturgeonStore) GetSeasons() ([]models.Season, error) {
	rows, err := s.db.Query(seasonsSql)

	seasons := []models.Season{}
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

var insertStructureSql = `insert into upload_site (site_id, site_fid, site_year, fieldoffice_id, 
	field_office, project_id, project, 
	segment_id, segment, season_id, season, bend, bendrn, bend_river_mile, comments,
	last_updated, upload_session_id, uploaded_by, upload_filename)
	values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19)`

func (s *PallidSturgeonStore) SaveSiteUpload(uploadSite models.UploadSite) error {
	_, err := s.db.Exec(insertStructureSql,
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

func (s *PallidSturgeonStore) SiteUploadSP(uploadedBy string, uploadSessionId int) error {
	uploadSiteStmt, err := s.db.Prepare("begin DATA_UPLOAD.uploadSiteDatasheetCheck (:1,:2);  end;")

	//var retVal string
	uploadSiteStmt.Exec(godror.PlSQLArrays, uploadedBy, uploadSessionId)

	//fmt.Println(retVal)

	return err
}
