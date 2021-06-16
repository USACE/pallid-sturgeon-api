package models

import "time"

type Season struct {
	ID           int    `db:"id" json:"id"`
	Code         string `db:"code" json:"code"`
	Description  string `db:"description" json:"description"`
	FieldAppFlag string `db:"field_app_flag" json:"fieldAppFlag"`
	ProjectCode  *int   `db:"project_code" json:"projectCode"`
}

type UploadSite struct {
	SiteID          int       `db:"site_id" json:"siteId"`
	SiteFID         string    `db:"site_fid" json:"siteFid"`
	SiteYear        int       `db:"site_year" json:"siteYear"`
	FieldofficeID   string    `db:"fieldoffice_id" json:"fieldofficeId"`
	FieldOffice     string    `db:"field_office" json:"fieldOffice"`
	ProjectId       int       `db:"project_id" json:"projectId"`
	Project         string    `db:"project" json:"project"`
	SegmentId       int       `db:"segment_id" json:"segmentId"`
	Segment         string    `db:"segment" json:"segment"`
	SeasonId        string    `db:"season_id" json:"seasonId"`
	Season          string    `db:"season" json:"season"`
	Bend            int       `db:"bend" json:"bend"`
	Bendrn          string    `db:"bendrn" json:"bendrn"`
	BendRiverMile   float64   `db:"bend_river_mile" json:"bendRiverMile"`
	Comments        string    `db:"comments" json:"comments"`
	LastUpdated     time.Time `db:"last_updated" json:"lastUpdated"`
	UploadSessionId int       `db:"upload_session_id" json:"uploadSessionId"`
	UploadedBy      string    `db:"uploaded_by" json:"uploadedBy"`
	UploadFilename  string    `db:"upload_filename" json:"uploadFilename"`
}
