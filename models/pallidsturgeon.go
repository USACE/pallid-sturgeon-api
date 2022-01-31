package models

import "time"

type Role struct {
	ID          int    `db:"id" json:"id"`
	Description string `db:"description" json:"description"`
}

type Season struct {
	ID           int    `db:"s_id" json:"id"`
	Code         string `db:"season_code" json:"code"`
	Description  string `db:"season_description" json:"description"`
	FieldAppFlag string `db:"field_app" json:"fieldAppFlag"`
	ProjectCode  *int   `db:"PROJECT_CODE" json:"projectCode"`
}

type FieldOffice struct {
	ID          int    `db:"FO_ID" json:"id"`
	Code        string `db:"FIELD_OFFICE_CODE" json:"code"`
	Description string `db:"FIELD_OFFICE_DESCRIPTION" json:"description"`
	State       string `db:"state" json:"state"`
}

type SampleMethod struct {
	Code        string `db:"SAMPLE_TYPE_CODE" json:"code"`
	Description string `db:"SAMPLE_TYPE_DESCRIPTION" json:"description"`
}

type SampleUnitType struct {
	Code        string `db:"SAMPLE_UNIT_TYPE_CODE" json:"code"`
	Description string `db:"SAMPLE_UNIT_TYPE_DESCRIPTION" json:"description"`
}

type Segment struct {
	ID             int     `db:"s_id" json:"id"`
	Code           int     `db:"segment_code" json:"code"`
	Description    *string `db:"segment_description" json:"description"`
	Type           string  `db:"segment_type" json:"type"`
	RiverCode      int     `db:"river" json:"riverCode"`
	UpperRiverMile *string `db:"upper_river_mile" json:"upperRiverMile"`
	LowerRiverMile *string `db:"lower_river_mile" json:"lowerRiverMile"`
	Rpma           *int    `db:"rpma" json:"rpma"`
}

type Bend struct {
	ID             int     `db:"BRM_ID" json:"id"`
	BendNumber     int     `db:"BEND_NUM" json:"bendNumber"`
	Description    *string `db:"B_DESC" json:"description"`
	SegmentCode    int     `db:"B_SEGMENT" json:"segmentCode"`
	UpperRiverMile *string `db:"upper_river_mile" json:"upperRiverMile"`
	LowerRiverMile *string `db:"lower_river_mile" json:"lowerRiverMile"`
	State          string  `db:"state" json:"state"`
}

type Project struct {
	Code        int    `db:"project_code" json:"code"`
	Description string `db:"project_description" json:"description"`
}

type FishSummaryWithCount struct {
	Items      []FishSummary `json:"items"`
	TotalCount int           `json:"totalCount"`
}

type FishSummary struct {
	UniqueID        int     `db:"mr_id" json:"uniqueID"`
	FishID          int     `db:"f_id" json:"fishId"`
	Year            int     `db:"year" json:"year"`
	FieldOffice     string  `db:"FIELD_OFFICE_CODE" json:"fieldOffice"`
	Project         int     `db:"PROJECT_CODE" json:"project"`
	Segment         int     `db:"SEGMENT_CODE" json:"segment"`
	Season          string  `db:"SEASON_CODE" json:"season"`
	Bend            int     `db:"BEND_NUMBER" json:"bend"`
	Bendrn          string  `db:"BEND_R_OR_N" json:"bendrn"`
	BendRiverMile   float64 `db:"bend_river_mile" json:"bendRiverMile"`
	Panelhook       string  `db:"panelhook" json:"panelhook"`
	Species         string  `db:"SPECIES_CODE" json:"species"`
	HatcheryOrigin  string  `db:"HATCHERY_ORIGIN_CODE" json:"hatcheryOrigin"`
	CheckedBy       string  `db:"checkby" json:"checkedby"`
	EditInitials    string  `db:"edit_initials" json:"editInitials"`
	LastEditComment string  `db:"last_edit_comment" json:"lastEditComment"`
}

type SuppSummaryWithCount struct {
	Items      []SuppSummary `json:"items"`
	TotalCount int           `json:"totalCount"`
}

type SuppSummary struct {
	FishCode        string  `db:"fish_code" json:"fishCode"`
	UniqueID        int     `db:"mr_id" json:"uniqueID"`
	FishID          int     `db:"f_id" json:"fishId"`
	Year            int     `db:"year" json:"year"`
	SuppID          int     `db:"sid_display" json:"suppId"`
	FieldOffice     string  `db:"FIELD_OFFICE_CODE" json:"fieldOffice"`
	Project         int     `db:"PROJECT_CODE" json:"project"`
	Segment         int     `db:"SEGMENT_CODE" json:"segment"`
	Season          string  `db:"SEASON_CODE" json:"season"`
	Bend            int     `db:"BEND_NUMBER" json:"bend"`
	Bendrn          string  `db:"BEND_R_OR_N" json:"bendrn"`
	BendRiverMile   float64 `db:"bend_river_mile" json:"bendRiverMile"`
	HatcheryOrigin  string  `db:"HATCHERY_ORIGIN_CODE" json:"hatcheryOrigin"`
	CheckedBy       string  `db:"checkby" json:"checkedby"`
	EditInitials    string  `db:"edit_initials" json:"editInitials"`
	LastEditComment string  `db:"last_edit_comment" json:"lastEditComment"`
}

type MissouriSummaryWithCount struct {
	Items      []MissouriSummary `json:"items"`
	TotalCount int               `json:"totalCount"`
}

type MissouriSummary struct {
	UniqueID        int       `db:"mr_id" json:"uniqueID"`
	Year            int       `db:"year" json:"year"`
	FieldOffice     string    `db:"FIELD_OFFICE_CODE" json:"fieldOffice"`
	Project         int       `db:"PROJECT_CODE" json:"project"`
	Segment         int       `db:"SEGMENT_CODE" json:"segment"`
	Season          string    `db:"SEASON_CODE" json:"season"`
	Bend            int       `db:"BEND_NUMBER" json:"bend"`
	Bendrn          string    `db:"BEND_R_OR_N" json:"bendrn"`
	BendRiverMile   float64   `db:"bend_river_mile" json:"bendRiverMile"`
	Subsample       int       `db:"subsample" json:"subsample"`
	Pass            int       `db:"subsample_pass" json:"pass"`
	SetDate         time.Time `db:"set_date" json:"setDate"`
	Conductivity    *string   `db:"conductivity" json:"conductivity"`
	CheckedBy       string    `db:"checkby" json:"checkedby"`
	EditInitials    string    `db:"edit_initials" json:"editInitials"`
	LastEditComment string    `db:"last_edit_comment" json:"lastEditComment"`
}

type GeneticSummaryWithCount struct {
	Items      []GeneticSummary `json:"items"`
	TotalCount int              `json:"totalCount"`
}

type GeneticSummary struct {
	Year               int       `db:"year" json:"year"`
	FieldOffice        string    `db:"FIELD_OFFICE_CODE" json:"fieldOffice"`
	Project            int       `db:"PROJECT_CODE" json:"project"`
	SturgeonType       string    `db:"sturgeon_type" json:"sturgeonType"`
	GeneticsVialNumber string    `db:"genetics_vial_number" json:"GeneticsVialNumber"`
	PitTag             string    `db:"pit_tag" json:"pitTag"`
	River              string    `db:"river" json:"river"`
	RiverMile          float64   `db:"river_mile" json:"riverMile"`
	State              string    `db:"state" json:"state"`
	SetDate            time.Time `db:"set_date" json:"setDate"`
	Broodstock         string    `db:"broodstock_yn" json:"broodstock"`
	HatchWild          string    `db:"hatchwild_yn" json:"hatchWild"`
	SpeciesID          string    `db:"speciesid_yn" json:"speciesId"`
	Archive            string    `db:"archive_yn" json:"archive"`
	EditInitials       string    `db:"edit_initials" json:"editInitials"`
	LastEditComment    string    `db:"last_edit_comment" json:"lastEditComment"`
}

type SearchSummaryWithCount struct {
	Items      []SearchSummary `json:"items"`
	TotalCount int             `json:"totalCount"`
}

type SearchSummary struct {
	SeID            int     `db:"se_id" json:"seId"`
	SearchDate      string  `db:"Search_date" json:"searchDate"`
	Recorder        string  `db:"recorder" json:"recorder"`
	SearchTypeCode  string  `db:"search_type_code" json:"searchTypeCode"`
	StartTime       string  `db:"start_time" json:"startTime"`
	StartLatitude   float64 `db:"start_latitude" json:"startLatitude"`
	StartLongitude  float64 `db:"start_longitude" json:"startLongitude"`
	StopTime        string  `db:"stop_time" json:"stopTime"`
	StopLatitude    float64 `db:"stop_latitude" json:"stopLatitude"`
	StopLongitude   float64 `db:"stop_longitude" json:"stopLongitude"`
	SeFID           string  `db:"se_fid" json:"seFid"`
	DsID            int     `db:"ds_id" json:"dsId"`
	SiteFID         string  `db:"site_fid" json:"siteFid"`
	Temp            *string `db:"temp" json:"temp"`
	Conductivity    *string `db:"conductivity" json:"conductivity"`
	EditInitials    string  `db:"edit_initials" json:"editInitials"`
	LastEditComment string  `db:"last_edit_comment" json:"lastEditComment"`
}

type SummaryWithCount struct {
	Items      []map[string]string `json:"items"`
	TotalCount int                 `json:"totalCount"`
}

type Upload struct {
	EditInitials       string                 `db:"edit_initials" json:"editInitials"`
	SiteUpload         UploadSiteData         `json:"siteUpload"`
	FishUpload         UploadFishData         `json:"fishUpload"`
	SearchUpload       UploadSearchData       `json:"searchUpload"`
	ProcedureUpload    UploadProcedureData    `json:"procedureUpload"`
	UploadSupplemental UploadSupplementalData `json:"uploadSupplemental"`
	MoriverUpload      UploadMoriverData      `json:"moriverUpload"`
	TelemetryUpload    UploadTelemetryData    `json:"telemetryUpload"`
}

type SiteDataEntryWithCount struct {
	Items      []UploadSite `json:"items"`
	TotalCount int          `json:"totalCount"`
}

type UploadSiteData struct {
	Items          []UploadSite `json:"items"`
	UploadFilename string       `db:"upload_filename" json:"uploadFilename"`
}

type UploadSite struct {
	BrmID              int       `db:"brm_id" json:"brmId"`
	SiteID             int       `db:"site_id" json:"siteId"`
	SiteFID            string    `db:"site_fid" json:"siteFid"`
	SiteYear           int       `db:"site_year" json:"siteYear"`
	FieldofficeID      string    `db:"fieldoffice_id" json:"fieldofficeId"`
	FieldOffice        string    `db:"field_office" json:"fieldOffice"`
	ProjectId          int       `db:"project_id" json:"projectId"`
	Project            string    `db:"project" json:"project"`
	SegmentId          int       `db:"segment_id" json:"segmentId"`
	Segment            string    `db:"segment" json:"segment"`
	SeasonId           string    `db:"season_id" json:"seasonId"`
	Season             string    `db:"season" json:"season"`
	SampleUnitTypeCode string    `db:"SAMPLE_UNIT_TYPE" json:"sampleUnitTypeCode"`
	Bend               int       `db:"bend" json:"bend"`
	Bendrn             string    `db:"bendrn" json:"bendrn"`
	BendRiverMile      float64   `db:"bend_river_mile" json:"bendRiverMile"`
	EditInitials       string    `db:"edit_initials" json:"editInitials"`
	LastUpdated        time.Time `db:"last_updated" json:"lastUpdated"`
	UploadSessionId    int       `db:"upload_session_id" json:"uploadSessionId"`
	UploadedBy         string    `db:"uploaded_by" json:"uploadedBy"`
	UploadFilename     string    `db:"upload_filename" json:"uploadFilename"`
}

type SitesWithCount struct {
	Items      []Sites `json:"items"`
	TotalCount int     `json:"totalCount"`
}

type Sites struct {
	SiteID             int       `db:"site_id" json:"siteId"`
	SiteFID            string    `db:"site_fid" json:"siteFid"`
	SiteYear           int       `db:"site_year" json:"siteYear"`
	FieldofficeID      string    `db:"fieldoffice_id" json:"fieldofficeId"`
	FieldOffice        string    `db:"field_office" json:"fieldOffice"`
	ProjectId          int       `db:"project_id" json:"projectId"`
	Project            string    `db:"project" json:"project"`
	SegmentId          int       `db:"segment_id" json:"segmentId"`
	Segment            string    `db:"segment" json:"segment"`
	SeasonId           string    `db:"season_id" json:"seasonId"`
	Season             string    `db:"season" json:"season"`
	SampleUnitTypeCode string    `db:"SAMPLE_UNIT_TYPE" json:"sampleUnitTypeCode"`
	Bend               int       `db:"bend" json:"bend"`
	Bendrn             string    `db:"bendrn" json:"bendrn"`
	BendRiverMile      float64   `db:"brm_id" json:"bendRiverMile"`
	EditInitials       string    `db:"edit_initials" json:"editInitials"`
	LastUpdated        time.Time `db:"last_updated" json:"lastUpdated"`
	UploadSessionId    int       `db:"upload_session_id" json:"uploadSessionId"`
	UploadedBy         string    `db:"uploaded_by" json:"uploadedBy"`
	UploadFilename     string    `db:"upload_filename" json:"uploadFilename"`
}

type FishDataEntryWithCount struct {
	Items      []UploadFish `json:"items"`
	TotalCount int          `json:"totalCount"`
}

type UploadFishData struct {
	Items          []UploadFish `json:"items"`
	UploadFilename string       `db:"upload_filename" json:"uploadFilename"`
}

// type FishDataEntry struct {
// 	Fid         int      `db:"f_id" json:"fid"`
// 	Ffid        string   `db:"f_fid" json:"ffid"`
// 	Fieldoffice string   `db:"FIELDOFFICE" json:"fieldOffice"`
// 	Project     *int     `db:"PROJECT_ID" json:"project"`
// 	Segment     *int     `db:"SEGMENT_ID" json:"segment"`
// 	UniqueID    *int     `db:"uniqueidentifier" json:"uniqueID"`
// 	Id          *int     `db:"id" json:"id"`
// 	Panelhook   string   `db:"panelhook" json:"panelhook"`
// 	Bait        string   `db:"bait" json:"bait"`
// 	Species     string   `db:"SPECIES_ID" json:"species"`
// 	Length      *float32 `db:"length" json:"length"`
// 	Weight      *float32 `db:"weight" json:"weight"`
// 	Fishcount   *int     `db:"FISHCOUNT" json:"fishcount"`
// 	Otolith     string   `db:"otolith" json:"otolith"`
// 	Rayspine    string   `db:"rayspine" json:"rayspine"`
// 	Scale       string   `db:"scale" json:"scale"`
// 	Ftprefix    string   `db:"FTPREFIX" json:"ftprefix"`
// 	Ftnum       string   `db:"FTNUM" json:"ftnum"`
// 	Ftmr        string   `db:"FTMR" json:"ftmr"`
// 	MrID        *int     `db:"mr_id" json:"mrId"`
// }

type UploadFish struct {
	Id              *int      `db:"id" json:"id"`
	SiteID          int       `db:"site_id" json:"siteId"`
	MrFid           string    `db:"mr_fid" json:"mrFid"`
	Fid             int       `db:"f_id" json:"fid"`
	Ffid            string    `db:"f_fid" json:"ffid"`
	MrID            *int      `db:"mr_id" json:"mrId"`
	Panelhook       string    `db:"panelhook" json:"panelhook"`
	Bait            string    `db:"bait" json:"bait"`
	Species         string    `db:"species" json:"species"`
	Length          *float32  `db:"length" json:"length"`
	Weight          *float32  `db:"weight" json:"weight"`
	Fishcount       int       `db:"fishcount" json:"fishcount"`
	FinCurl         string    `db:"fin_curl" json:"finCurl"`
	Otolith         string    `db:"otolith" json:"otolith"`
	Rayspine        string    `db:"rayspine" json:"rayspine"`
	Scale           string    `db:"scale" json:"scale"`
	Ftprefix        string    `db:"ftprefix" json:"ftprefix"`
	Ftnum           string    `db:"ftnum" json:"ftnum"`
	Ftmr            string    `db:"ftmr" json:"ftmr"`
	Comments        string    `db:"comments" json:"comments"`
	LastUpdated     time.Time `db:"last_updated" json:"lastUpdated"`
	UploadSessionId int       `db:"upload_session_id" json:"uploadSessionId"`
	EditInitials    string    `db:"edit_initials" json:"editInitials"`
	LastEditComment string    `db:"last_edit_comment" json:"lastEditComment"`
	UploadedBy      string    `db:"uploaded_by" json:"uploadedBy"`
	UploadFilename  string    `db:"upload_filename" json:"uploadFilename"`
	Project         *int      `db:"PROJECT_ID" json:"project"`
	UniqueID        *int      `db:"uniqueidentifier" json:"uniqueID"`
	Segment         *int      `db:"SEGMENT_ID" json:"segment"`
	Fieldoffice     string    `db:"FIELDOFFICE" json:"fieldOffice"`
}

type UploadSearchData struct {
	Items          []UploadSearch `json:"items"`
	UploadFilename string         `db:"upload_filename" json:"uploadFilename"`
}

type UploadSearch struct {
	SiteID          int       `db:"site_id" json:"siteId"`
	SeFid           string    `db:"se_fid" json:"seFid"`
	DsId            int       `db:"ds_id" json:"dsId"`
	SiteFid         string    `db:"site_fid" json:"siteFid"`
	SearchDate      time.Time `db:"search_date" json:"searchDate"`
	Recorder        string    `db:"recorder" json:"recorder"`
	SearchTypeCode  string    `db:"search_type_code" json:"searchTypeCode"`
	SearchDay       int       `db:"search_day" json:"searchDay"`
	StartTime       string    `db:"start_time" json:"startTime"`
	StartLatitude   float64   `db:"start_latitude" json:"startLatitude"`
	StartLongitude  float64   `db:"start_longitude" json:"startLongitude"`
	StopTime        string    `db:"stop_time" json:"stopTime"`
	StopLatitude    float64   `db:"stop_latitude" json:"stopLatitude"`
	StopLongitude   float64   `db:"stop_longitude" json:"stopLongitude"`
	Temp            string    `db:"temp" json:"temp"`
	Conductivity    float64   `db:"conductivity" json:"conductivity"`
	EditInitials    string    `db:"edit_initials" json:"editInitials"`
	LastEditComment string    `db:"last_edit_comment" json:"lastEditComment"`
	LastUpdated     time.Time `db:"last_updated" json:"lastUpdated"`
	UploadSessionId int       `db:"upload_session_id" json:"uploadSessionId"`
	UploadedBy      string    `db:"uploaded_by" json:"uploadedBy"`
	UploadFilename  string    `db:"upload_filename" json:"uploadFilename"`
}

type UploadProcedureData struct {
	Items          []UploadProcedure `json:"items"`
	UploadFilename string            `db:"upload_filename" json:"uploadFilename"`
}

type UploadProcedure struct {
	Id                        int       `db:"id" json:"id"`
	FFid                      string    `db:"f_fid" json:"f_fid"`
	PurposeCode               string    `db:"PURPOSE" json:"purposeCode"`
	ProcedurDate              time.Time `db:"SURGERY_DATE" json:"procedurDate"`
	ProcedureStartTime        string    `db:"procedure_start_time" json:"procedureStartTime"`
	ProcedureEndTime          string    `db:"procedure_end_time" json:"procedureEndTime"`
	ProcedureBy               string    `db:"procedure_by" json:"procedureBy"`
	AntibioticInjectionInd    int       `db:"ANTIBIOTIC_INJECTION" json:"antibioticInjectionInd"`
	PhotoDorsalInd            int       `db:"P_DORSAL" json:"photoDorsalInd"`
	PhotoVentralInd           int       `db:"P_VENTRAL" json:"photoVentralInd"`
	PhotoLeftInd              int       `db:"P_LEFT" json:"photoLeftInd"`
	OldRadioTagNum            int       `db:"old_radio_tag_num" json:"oldRadioTagNum"`
	OldFrequencyId            int       `db:"OLD_FREQ_ID" json:"oldFrequencyId"`
	DstSerialNum              int       `db:"dst_serial_num" json:"dstSerialNum"`
	DstStartDate              time.Time `db:"dst_start_date" json:"dstStartDate"`
	DstStartTime              string    `db:"dst_start_time" json:"dstStartTime"`
	DstReimplantInd           int       `db:"DST_REIMPLANT" json:"dstReimplantInd"`
	NewRadioTagNum            int       `db:"new_radio_tag_num" json:"newRadioTagNum"`
	NewFrequencyId            int       `db:"NEW_FREQ_ID" json:"newFrequencyId"`
	SexCode                   string    `db:"SEX" json:"sexCode"`
	BloodSampleInd            int       `db:"BLOOD_SAMPLE" json:"bloodSampleInd"`
	EggSampleInd              int       `db:"EGG_SAMPLE" json:"eggSampleInd"`
	Comments                  string    `db:"comments" json:"comments"`
	FishHealthComments        string    `db:"FISH_HEALTH_COMMENT" json:"fishHealthComments"`
	EvalLocationCode          string    `db:"EVAL_LOCATION" json:"evalLocationCode"`
	SpawnCode                 string    `db:"SPAWN_STATU" json:"spawnCode"`
	VisualReproStatusCode     string    `db:"VISUAL_REPRO_STATUS" json:"visualReproStatusCode"`
	UltrasoundReproStatusCode string    `db:"ULTRASOUND_REPRO_STATUS" json:"ultrasoundReproStatusCode"`
	ExpectedSpawnYear         int       `db:"expected_spawn_year" json:"expectedSpawnYear"`
	UltrasoundGonadLength     float64   `db:"ultrasound_gonad_length" json:"ultrasoundGonadLength"`
	GonadCondition            string    `db:"gonad_condition" json:"gonadCondition"`
	EditInitials              string    `db:"edit_initials" json:"editInitials"`
	LastEditComment           string    `db:"last_edit_comment" json:"lastEditComment"`
	LastUpdated               time.Time `db:"last_updated" json:"lastUpdated"`
	UploadSessionId           int       `db:"upload_session_id" json:"uploadSessionId"`
	UploadedBy                string    `db:"uploaded_by" json:"uploadedBy"`
	UploadFilename            string    `db:"upload_filename" json:"uploadFilename"`
}

type SupplementalDataEntryWithCount struct {
	Items      []UploadSupplemental `json:"items"`
	TotalCount int                  `json:"totalCount"`
}

type UploadSupplementalData struct {
	Items          []UploadSupplemental `json:"items"`
	UploadFilename string               `db:"upload_filename" json:"uploadFilename"`
}

type UploadSupplemental struct {
	Id                 *int      `db:"id" json:"id"`
	Sid                int       `db:"s_id" json:"sid"`
	Fid                int       `db:"f_id" json:"fid"`
	SiteID             int       `db:"site_id" json:"siteId"`
	FFid               string    `db:"f_fid" json:"fFid"`
	MrId               string    `db:"mr_id" json:"mrId"`
	MrFid              string    `db:"mr_fid" json:"mrFid"`
	Tagnumber          string    `db:"tagnumber" json:"tagnumber"`
	Pitrn              string    `db:"pitrn" json:"pitrn"`
	Scuteloc           string    `db:"scuteloc" json:"scuteloc"`
	Scutenum           *string   `db:"scutenum" json:"scutenum"`
	Scuteloc2          string    `db:"scuteloc2" json:"scuteloc2"`
	Scutenum2          *string   `db:"scutenum2" json:"scutenum2"`
	Elhv               string    `db:"elhv" json:"elhv"`
	Elcolor            string    `db:"elcolor" json:"elcolor"`
	Erhv               string    `db:"erhv" json:"erhv"`
	Ercolor            string    `db:"ercolor" json:"ercolor"`
	Cwtyn              string    `db:"cwtyn" json:"cwtyn"`
	Dangler            string    `db:"dangler" json:"dangler"`
	Genetic            string    `db:"genetic" json:"genetic"`
	GeneticsVialNumber string    `db:"genetics_vial_number" json:"geneticsVialNumber"`
	Broodstock         *string   `db:"broodstock" json:"broodstock"`
	HatchWild          *string   `db:"hatch_wild" json:"hatchWild"`
	SpeciesId          *int      `db:"species_id" json:"speciesId"`
	Head               *string   `db:"head" json:"head"`
	Snouttomouth       *string   `db:"snouttomouth" json:"snouttomouth"`
	Inter              *string   `db:"inter" json:"inter"`
	Mouthwidth         *string   `db:"mouthwidth" json:"mouthwidth"`
	MIb                *string   `db:"m_ib" json:"mIb"`
	LOb                *string   `db:"l_ob" json:"lOb"`
	LIb                *string   `db:"l_ib" json:"lIb"`
	RIb                *string   `db:"r_ib" json:"rIb"`
	ROb                *string   `db:"r_ob" json:"rOb"`
	Anal               *string   `db:"anal" json:"anal"`
	Dorsal             *string   `db:"dorsal" json:"dorsal"`
	Status             string    `db:"status" json:"status"`
	HatcheryOrigin     string    `db:"hatchery_origin" json:"hatcheryOrigin"`
	Sex                string    `db:"sex" json:"sex"`
	Stage              string    `db:"stage" json:"stage"`
	Recapture          string    `db:"recapture" json:"recapture"`
	Photo              string    `db:"photo" json:"photo"`
	GeneticNeeds       string    `db:"genetic_needs" json:"geneticNeeds"`
	OtherTagInfo       string    `db:"other_tag_info" json:"otherTagInfo"`
	Comments           string    `db:"comments" json:"comments"`
	EditInitials       string    `db:"edit_initials" json:"editInitials"`
	LastEditComment    string    `db:"last_edit_comment" json:"lastEditComment"`
	LastUpdated        time.Time `db:"last_updated" json:"lastUpdated"`
	UploadSessionId    int       `db:"upload_session_id" json:"uploadSessionId"`
	UploadedBy         string    `db:"uploaded_by" json:"uploadedBy"`
	UploadFilename     string    `db:"upload_filename" json:"uploadFilename"`
}

type MoriverDataEntryWithCount struct {
	Items      []UploadMoriver `json:"items"`
	TotalCount int             `json:"totalCount"`
}

type UploadMoriverData struct {
	Items          []UploadMoriver `json:"items"`
	UploadFilename string          `db:"upload_filename" json:"uploadFilename"`
}

type UploadMoriver struct {
	SiteID           int       `db:"site_id" json:"siteId"`
	SiteFid          string    `db:"site_fid" json:"siteFid"`
	MrID             int       `db:"mr_id" json:"mrId"`
	MrFid            string    `db:"mr_fid" json:"mrFid"`
	Season           string    `db:"season" json:"season"`
	Setdate          time.Time `db:"setdate" json:"setdate"`
	Subsample        float64   `db:"subsample" json:"subsample"`
	Subsamplepass    float64   `db:"subsamplepass" json:"subsamplepass"`
	SubsampleROrN    string    `db:"subsample_r_or_n" json:"subsampleROrN"`
	Subsamplen       string    `db:"subsamplen" json:"subsamplen"`
	Recorder         string    `db:"recorder" json:"recorder"`
	Gear             string    `db:"gear" json:"gear"`
	GearType         string    `db:"gear_type" json:"gearType"`
	Temp             *string   `db:"temp" json:"temp"`
	Turbidity        *string   `db:"turbidity" json:"turbidity"`
	Conductivity     *string   `db:"conductivity" json:"conductivity"`
	Do               *string   `db:"do" json:"do"`
	Distance         *string   `db:"distance" json:"distance"`
	Width            *string   `db:"width" json:"width"`
	Netrivermile     *string   `db:"netrivermile" json:"netrivermile"`
	Structurenumber  string    `db:"structurenumber" json:"structurenumber"`
	Usgs             string    `db:"usgs" json:"usgs"`
	Riverstage       *string   `db:"riverstage" json:"riverstage"`
	Discharge        *string   `db:"discharge" json:"discharge"`
	U1               string    `db:"u1" json:"u1"`
	U2               string    `db:"u2" json:"u2"`
	U3               string    `db:"u3" json:"u3"`
	U4               string    `db:"u4" json:"u4"`
	U5               string    `db:"u5" json:"u5"`
	U6               string    `db:"u6" json:"u6"`
	U7               string    `db:"u7" json:"u7"`
	Macro            string    `db:"macro" json:"macro"`
	Meso             string    `db:"meso" json:"meso"`
	Habitatrn        string    `db:"habitatrn" json:"habitatrn"`
	Qc               string    `db:"qc" json:"qc"`
	MicroStructure   string    `db:"micro_structure" json:"microStructure"`
	StructureFlow    string    `db:"structure_flow" json:"structureFlow"`
	StructureMod     string    `db:"structure_mod" json:"structureMod"`
	SetSite1         string    `db:"set_site_1" json:"setSite_1"`
	SetSite2         string    `db:"set_site_2" json:"setSite_2"`
	SetSite3         string    `db:"set_site_3" json:"setSite_3"`
	StartTime        string    `db:"starttime" json:"startTime"`
	StartLatitude    *string   `db:"startlatitude" json:"startLatitude"`
	StartLongitude   *string   `db:"startlongitude" json:"startLongitude"`
	StopTime         string    `db:"stoptime" json:"stopTime"`
	StopLatitude     *string   `db:"stoplatitude" json:"stopLatitude"`
	StopLongitude    *string   `db:"stop_longitude" json:"stopLongitude"`
	Depth1           *string   `db:"depth1" json:"depth1"`
	Velocitybot1     *string   `db:"velocitybot1" json:"velocitybot1"`
	Velocity08_1     *string   `db:"velocity08_1" json:"velocity08_1"`
	Velocity02or06_1 *string   `db:"velocity02or06_1" json:"velocity02or06_1"`
	Depth2           *string   `db:"depth2" json:"depth2"`
	Velocitybot2     *string   `db:"velocitybot2" json:"velocitybot2"`
	Velocity08_2     *string   `db:"velocity08_2" json:"velocity08_2"`
	Velocity02or06_2 *string   `db:"velocity02or06_2" json:"velocity02or06_2"`
	Depth3           *string   `db:"depth3" json:"depth3"`
	Velocitybot3     *string   `db:"velocitybot3" json:"velocitybot3"`
	Velocity08_3     *string   `db:"velocity08_3" json:"velocity08_3"`
	Velocity02or06_3 *string   `db:"velocity02or06_3" json:"velocity02or06_3"`
	Watervel         *string   `db:"watervel" json:"watervel"`
	Cobble           *string   `db:"cobble" json:"cobble"`
	Organic          *string   `db:"organic" json:"organic"`
	Silt             *string   `db:"silt" json:"silt"`
	Sand             *string   `db:"sand" json:"sand"`
	Gravel           *string   `db:"gravel" json:"gravel"`
	Comments         string    `db:"comments" json:"comments"`
	Complete         *string   `db:"complete" json:"complete"`
	Checkby          string    `db:"checkby" json:"checkby"`
	NoTurbidity      string    `db:"no_turbidity" json:"noTurbidity"`
	NoVelocity       string    `db:"no_velocity" json:"noVelocity"`
	EditInitials     string    `db:"edit_initials" json:"editInitials"`
	LastEditComment  string    `db:"last_edit_comment" json:"lastEditComment"`
	Project          *int      `db:"PROJECT_ID" json:"project"`
	FieldOffice      string    `db:"FIELDOFFICE" json:"fieldOffice"`
	Segment          *int      `db:"SEGMENT_ID" json:"segment"`
	LastUpdated      time.Time `db:"last_updated" json:"lastUpdated"`
	UploadSessionId  int       `db:"upload_session_id" json:"uploadSessionId"`
	UploadedBy       string    `db:"uploaded_by" json:"uploadedBy"`
	UploadFilename   string    `db:"upload_filename" json:"uploadFilename"`
}

type UploadTelemetryData struct {
	Items          []UploadTelemetry `json:"items"`
	UploadFilename string            `db:"upload_filename" json:"uploadFilename"`
}

type UploadTelemetry struct {
	TFid               string    `db:"t_fid" json:"tFid"`
	SeFid              string    `db:"se_fid" json:"seFid"`
	Bend               float64   `db:"bend" json:"bend"`
	RadioTagNum        int       `db:"radio_tag_num" json:"radioTagNum"`
	FrequencyIdCode    int       `db:"frequency_id_code" json:"frequencyIdCode"`
	CaptureTime        string    `db:"capture_time" json:"captureTime"`
	CaptureLatitude    float64   `db:"capture_latitude" json:"captureLatitude"`
	CaptureLongitude   float64   `db:"capture_longitude" json:"captureLongitude"`
	PositionConfidence float64   `db:"position_confidence" json:"positionConfidence"`
	MacroId            string    `db:"macro_id" json:"macroId"`
	MesoId             string    `db:"meso_id" json:"mesoId"`
	Depth              float64   `db:"depth" json:"depth"`
	Temp               float64   `db:"temp" json:"temp"`
	Conductivity       float64   `db:"conductivity" json:"conductivity"`
	Turbidity          float64   `db:"turbidity" json:"turbidity"`
	Silt               float64   `db:"silt" json:"silt"`
	Sand               float64   `db:"sand" json:"sand"`
	Gravel             float64   `db:"gravel" json:"gravel"`
	EditInitials       string    `db:"edit_initials" json:"editInitials"`
	LastEditComment    string    `db:"last_edit_comment" json:"lastEditComment"`
	Comments           string    `db:"comments" json:"comments"`
	LastUpdated        time.Time `db:"last_updated" json:"lastUpdated"`
	UploadSessionId    int       `db:"upload_session_id" json:"uploadSessionId"`
	UploadedBy         string    `db:"uploaded_by" json:"uploadedBy"`
	UploadFilename     string    `db:"upload_filename" json:"uploadFilename"`
}

type ProcedureOut struct {
	UploadSessionId   int    `json:"uploadSessionId"`
	UploadedBy        string `json:"uploadedBy"`
	SiteCntFinal      int    `json:"siteCntFinal"`
	MrCntFinal        int    `json:"mrCntFinal"`
	FishCntFinal      int    `json:"fishCntFinal"`
	SearchCntFinal    int    `json:"searchCntFinal"`
	SuppCntFinal      int    `json:"suppCntFinal"`
	TelemetryCntFinal int    `json:"telemetryCntFinal"`
	ProcedureCntFinal int    `json:"procedureCntFinal"`
	NoSiteCnt         int    `json:"noSiteCnt"`
	SiteMatch         int    `json:"siteMatch"`
	NoSiteIDMsg       string `json:"noSiteIDMsg"`
}

type UsgNoVialNumber struct {
	Fp                 string `db:"fp" json:"fp"`
	SpeciesCode        string `db:"SPECIES" json:"speciesCode"`
	FID                int    `db:"f_id" json:"fId"`
	MrID               int    `db:"mr_id" json:"mrID"`
	MrsiteID           int    `db:"mrsite_id" json:"mrsiteId"`
	SSiteID            int    `db:"s_site_id" json:"sSiteID"`
	FFID               string `db:"f_fid" json:"fFId"`
	GeneticsVialNumber string `db:"genetics_vial_number" json:"GeneticsVialNumber"`
}

type ErrorCount struct {
	Year  int `db:"year" json:"year"`
	Count int `db:"count(el.el_id)" json:"count"`
}

type DownloadInfo struct {
	Name        string `db:"name" json:"name"`
	DisplayName string `db:"display_name" json:"displayName"`
	LastUpdated string `db:"last_updated" json:"lastUpdated"`
}
