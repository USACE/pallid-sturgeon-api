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

type BendRn struct {
	ID          int    `db:"bs_id" json:"id"`
	Code        string `db:"bend_selection_code" json:"code"`
	Description string `db:"bend_selection_description" json:"description"`
}

type Project struct {
	Code        int    `db:"project_code" json:"code"`
	Description string `db:"project_description" json:"description"`
}

type Meso struct {
	Code string `db:"mesohabitat_code" json:"code"`
}

type StructureFlow struct {
	ID   int    `db:"structure_flow_code" json:"id"`
	Code string `db:"structure_flow" json:"code"`
}

type StructureMod struct {
	Code        string `db:"structure_mod_code" json:"code"`
	Description string `db:"structure_mod" json:"description"`
}

type Species struct {
	Code string `db:"alpha_code" json:"code"`
}

type FtPrefix struct {
	Code string `db:"tag_prefix_code" json:"code"`
}

type Mr struct {
	Code        string `db:"mark_recapture_code" json:"code"`
	Description string `db:"mark_recapture_description" json:"description"`
}

type Otolith struct {
	Code        string `db:"code" json:"code"`
	Description string `db:"description" json:"description"`
}

type HeaderData struct {
	SiteId         int     `db:"site_id" json:"siteId"`
	Year           int     `db:"year" json:"year"`
	FieldOffice    string  `db:"fieldoffice" json:"fieldoffice"`
	Project        int     `db:"project_id" json:"project"`
	Segment        int     `db:"segment_id" json:"segment"`
	Season         string  `db:"season" json:"season"`
	Bend           int     `db:"bend" json:"bend"`
	Bendrn         string  `db:"bendrn" json:"bendrn"`
	BendRiverMile  float64 `db:"bendrivermile" json:"bendrivermile"`
	SampleUnitType string  `db:"sample_unit_type" json:"sampleUnitType"`
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
	TagNumber       string  `db:"tag_number" json:"tagNumber"`
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

type ProcedureSummaryWithCount struct {
	Items      []ProcedureSummary `json:"items"`
	TotalCount int                `json:"totalCount"`
}

type ProcedureSummary struct {
	ID          int    `db:"pid_display" json:"id"`
	UniqueID    int    `db:"mr_id" json:"uniqueId"`
	Year        int    `db:"year" json:"year"`
	FieldOffice string `db:"field_office_code" json:"fieldOffice"`
	Project     int    `db:"project_code" json:"project"`
	Segment     int    `db:"segment_code" json:"segment"`
	Season      string `db:"season_code" json:"season"`
	PurposeCode string `db:"purpose_code" json:"purposeCode"`
	// ProcedureDate     time.Time `db:"procedure_date" json:"procedureDate"`
	NewRadioTagNum    int    `db:"new_radio_tag_num" json:"newRadioTagNum"`
	NewFrequencyId    int    `db:"new_frequency_id" json:"newFrequencyId"`
	SpawnCode         string `db:"spawn_code" json:"spawnCode"`
	ExpectedSpawnYear int    `db:"expected_spawn_year" json:"expectedSpawnYear"`
}

type TelemetrySummaryWithCount struct {
	Items      []TelemetrySummary `json:"items"`
	TotalCount int                `json:"totalCount"`
}

type TelemetrySummary struct {
	UniqueID           int     `db:"mr_id" json:"uniqueId"`
	TId                string  `db:"t_id" json:"tId"`
	TFid               string  `db:"t_fid" json:"tFid"`
	SeId               string  `db:"se_id" json:"seFid"`
	Year               int     `db:"year" json:"year"`
	FieldOffice        string  `db:"field_office_code" json:"fieldOffice"`
	Project            int     `db:"project_code" json:"project"`
	Segment            int     `db:"segment_code" json:"segment"`
	Season             string  `db:"season_code" json:"season"`
	Bend               float64 `db:"bend_number" json:"bend"`
	RadioTagNum        int     `db:"radio_tag_num" json:"radioTagNum"`
	FrequencyIdCode    int     `db:"frequency_id" json:"frequencyIdCode"`
	CaptureTime        string  `db:"capture_time" json:"captureTime"`
	CaptureLatitude    float64 `db:"capture_latitude" json:"captureLatitude"`
	CaptureLongitude   float64 `db:"capture_longitude" json:"captureLongitude"`
	PositionConfidence float64 `db:"position_confidence" json:"positionConfidence"`
	MacroId            string  `db:"macro_code" json:"macroId"`
	MesoId             string  `db:"meso_code" json:"mesoId"`
	Depth              float64 `db:"depth" json:"depth"`
	Conductivity       float64 `db:"conductivity" json:"conductivity"`
	Turbidity          float64 `db:"turbidity" json:"turbidity"`
}

type Upload struct {
	EditInitials       string                 `db:"edit_initials" json:"editInitials"`
	SiteUpload         UploadSiteData         `json:"siteUpload"`
	FishUpload         UploadFishData         `json:"fishUpload"`
	SearchUpload       UploadSearchData       `json:"searchUpload"`
	ProcedureUpload    UploadProcedureData    `json:"procedureUpload"`
	SupplementalUpload UploadSupplementalData `json:"supplementalUpload"`
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
	// BrmID              int       `db:"brm_id" json:"brmId"`
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
	SampleUnitTypeCode string    `db:"sample_unit_type" json:"sampleUnitTypeCode"`
	Bend               int       `db:"bend" json:"bend"`
	Bendrn             string    `db:"bendrn" json:"bendrn"`
	BendRiverMile      float64   `db:"bend_river_mile" json:"bendRiverMile"`
	Comments           string    `db:"comments" json:"comments"`
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
	Year               int       `db:"year" json:"year"`
	BrmID              int       `db:"brm_id" json:"brmId"`
	FieldofficeId      string    `db:"fieldoffice" json:"fieldoffice"`
	ProjectId          int       `db:"project_id" json:"projectId"`
	SegmentId          int       `db:"segment_id" json:"segmentId"`
	SeasonId           string    `db:"season" json:"season"`
	SampleUnitTypeCode string    `db:"sample_unit_type" json:"sampleUnitType"`
	Bend               int       `db:"bend" json:"bend"`
	Bendrn             string    `db:"bendrn" json:"bendrn"`
	BendRiverMile      *string   `db:"bend_river_mile" json:"bendRiverMile"`
	Complete           int       `db:"complete" json:"complete"`
	Approved           int       `db:"approved" json:"approved"`
	BkgColor           string    `db:"bkg_color" json:"bkgColor"`
	EditInitials       string    `db:"edit_initials" json:"editInitials"`
	LastEditComment    string    `db:"last_edit_comment" json:"last_edit_comment"`
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

type UploadFish struct {
	Id                 *int      `db:"id" json:"id"`
	SiteID             int       `db:"site_id" json:"siteId"`
	MrFid              string    `db:"mr_fid" json:"mrFid"`
	Fid                int       `db:"f_id" json:"fid"`
	Ffid               string    `db:"f_fid" json:"ffid"`
	MrID               *int      `db:"mr_id" json:"mrId"`
	Panelhook          string    `db:"panelhook" json:"panelHook"`
	Bait               string    `db:"bait" json:"bait"`
	Species            string    `db:"species" json:"species"`
	Length             *float32  `db:"length" json:"length"`
	Weight             *float32  `db:"weight" json:"weight"`
	Fishcount          int       `db:"fishcount" json:"countF"`
	FinCurl            string    `db:"fin_curl" json:"finCurl"`
	Otolith            string    `db:"otolith" json:"otolith"`
	Rayspine           string    `db:"rayspine" json:"raySpine"`
	Scale              string    `db:"scale" json:"scale"`
	Ftprefix           string    `db:"ftprefix" json:"ftPrefix"`
	Ftnum              string    `db:"ftnum" json:"ftnum"`
	Ftmr               string    `db:"ftmr" json:"mR"`
	Comments           string    `db:"comments" json:"comments"`
	Approved           int       `db:"approved" json:"approved"`
	LastUpdated        time.Time `db:"last_updated" json:"lastUpdated"`
	UploadSessionId    int       `db:"upload_session_id" json:"uploadSessionId"`
	EditInitials       string    `db:"edit_initials" json:"editInitials"`
	LastEditComment    string    `db:"last_edit_comment" json:"lastEditComment"`
	UploadedBy         string    `db:"uploaded_by" json:"uploadedBy"`
	UploadFilename     string    `db:"upload_filename" json:"uploadFilename"`
	Project            *int      `db:"PROJECT_ID" json:"project"`
	UniqueID           *int      `db:"uniqueidentifier" json:"uniqueID"`
	Segment            *int      `db:"SEGMENT_ID" json:"segment"`
	Fieldoffice        string    `db:"FIELDOFFICE" json:"fieldOffice"`
	GeneticsVialNumber string    `db:"genetics_vial_number" json:"geneticsVialNumber"`
}

type SearchDataEntryWithCount struct {
	Items      []UploadSearch `json:"items"`
	TotalCount int            `json:"totalCount"`
}

type UploadSearchData struct {
	Items          []UploadSearch `json:"items"`
	UploadFilename string         `db:"upload_filename" json:"uploadFilename"`
	TotalCount     int            `json:"totalCount"`
}

type UploadSearch struct {
	SiteId          int       `db:"site_id" json:"siteId"`
	SeId            int       `db:"se_id" json:"seId"`
	SeFid           string    `db:"se_fid" json:"seFid"`
	DsId            int       `db:"ds_id" json:"dsId"`
	SiteFid         string    `db:"site_fid" json:"siteFid"`
	SearchDate      string    `db:"search_date" json:"searchDate"`
	SearchDateTime  time.Time `db:"search_date" json:"searchDateTime"`
	Recorder        string    `db:"recorder" json:"recorder"`
	SearchTypeCode  string    `db:"search_type_code" json:"searchTypeCode"`
	SearchDay       int       `db:"search_day" json:"searchDay"`
	StartTime       string    `db:"start_time" json:"startTime"`
	StartLatitude   float64   `db:"start_latitude" json:"startLatitude"`
	StartLongitude  float64   `db:"start_longitude" json:"startLongitude"`
	StopTime        string    `db:"stop_time" json:"stopTime"`
	StopLatitude    float64   `db:"stop_latitude" json:"stopLatitude"`
	StopLongitude   float64   `db:"stop_longitude" json:"stopLongitude"`
	Temp            float64   `db:"temp" json:"temp"`
	Conductivity    float64   `db:"conductivity" json:"conductivity"`
	Checkby         string    `db:"checkby" json:"checkby"`
	EditInitials    string    `db:"edit_initials" json:"editInitials"`
	LastEditComment string    `db:"last_edit_comment" json:"lastEditComment"`
	LastUpdated     time.Time `db:"last_updated" json:"lastUpdated"`
	UploadSessionId int       `db:"upload_session_id" json:"uploadSessionId"`
	UploadedBy      string    `db:"uploaded_by" json:"uploadedBy"`
	UploadFilename  string    `db:"upload_filename" json:"uploadFilename"`
}

type ProcedureDataEntryWithCount struct {
	Items      []UploadProcedure `json:"items"`
	TotalCount int               `json:"totalCount"`
}

type UploadProcedureData struct {
	Items          []UploadProcedure `json:"items"`
	UploadFilename string            `db:"upload_filename" json:"uploadFilename"`
}

type UploadProcedure struct {
	Id                        int       `db:"id" json:"id"`
	Fid                       int       `db:"f_id" json:"fid"`
	FFid                      string    `db:"f_fid" json:"fFid"`
	MrFid                     string    `db:"MR_FID" json:"mrFid"`
	SiteID                    int       `db:"site_id" json:"siteId"`
	PurposeCode               string    `db:"purpose_code" json:"purpose"`
	ProcedureDate             string    `db:"PROCEDURE_DATE" json:"procedureDate"`
	ProcedureDateTime         time.Time `db:"PROCEDURE_DATE" json:"procedureDateTime"`
	ProcedureStartTime        string    `db:"procedure_start_time" json:"procedureStartTime"`
	ProcedureEndTime          string    `db:"procedure_end_time" json:"procedureEndTime"`
	ProcedureBy               string    `db:"procedure_by" json:"procedureBy"`
	AntibioticInjectionInd    int       `db:"ANTIBIOTIC_INJECTION_IND" json:"antibioticInjection"`
	PhotoDorsalInd            int       `db:"PHOTO_DORSAL_IND" json:"pDorsal"`
	PhotoVentralInd           int       `db:"PHOTO_VENTRAL_IND" json:"pVentral"`
	PhotoLeftInd              int       `db:"PHOTO_LEFT_IND" json:"pLeft"`
	OldRadioTagNum            int       `db:"old_radio_tag_num" json:"oldRadioTagNum"`
	OldFrequencyId            int       `db:"OLD_FREQUENCY_ID" json:"oldFrequencyId"`
	DstSerialNum              int       `db:"dst_serial_num" json:"dstSerialNum"`
	DstStartDate              string    `db:"dst_start_date" json:"dstStartDate"`
	DstStartDateTime          time.Time `db:"dst_start_date" json:"dstStartDateTime"`
	DstStartTime              string    `db:"dst_start_time" json:"dstStartTime"`
	DstReimplantInd           int       `db:"DST_REIMPLANT_IND" json:"dstReimplant"`
	NewRadioTagNum            int       `db:"new_radio_tag_num" json:"newRadioTagNum"`
	NewFrequencyId            int       `db:"NEW_FREQUENCY_ID" json:"newFreqId"`
	SexCode                   string    `db:"SEX_CODE" json:"sexCode"`
	BloodSampleInd            int       `db:"BLOOD_SAMPLE_IND" json:"bloodSample"`
	EggSampleInd              int       `db:"EGG_SAMPLE" json:"eggSample"`
	Comments                  string    `db:"comments" json:"comments"`
	FishHealthComments        string    `db:"FISH_HEALTH_COMMENTS" json:"fishHealthComment"`
	EvalLocationCode          string    `db:"EVAL_LOCATION_CODE" json:"evalLocation"`
	SpawnStatus               string    `db:"SPAWN_CODE" json:"spawnStatus"`
	VisualReproStatusCode     string    `db:"VISUAL_REPRO_STATUS" json:"visualReproStatus"`
	UltrasoundReproStatusCode string    `db:"ULTRASOUND_REPRO_STATUS" json:"ultrasoundReproStatus"`
	ExpectedSpawnYear         int       `db:"EXPECTED_SPAWN_YEAR" json:"expectedSpawnYear"`
	UltrasoundGonadLength     float64   `db:"ultrasound_gonad_length" json:"ultrasoundGonadLength"`
	GonadCondition            string    `db:"gonad_condition" json:"gonadCondition"`
	EditInitials              string    `db:"edit_initials" json:"editInitials"`
	LastEditComment           string    `db:"last_edit_comment" json:"lastEditComment"`
	LastUpdated               time.Time `db:"last_updated" json:"lastUpdated"`
	UploadSessionId           int       `db:"upload_session_id" json:"uploadSessionId"`
	UploadedBy                string    `db:"uploaded_by" json:"uploadedBy"`
	UploadFilename            string    `db:"upload_filename" json:"uploadFilename"`
	Checkby                   string    `db:"checkby" json:"checkby"`
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
	MrId               int       `db:"mr_id" json:"mrId"`
	MrFid              string    `db:"mr_fid" json:"mrFid"`
	Tagnumber          string    `db:"tagnumber" json:"tagnumber"`
	Pitrn              string    `db:"pitrn" json:"pitrn"`
	Scuteloc           string    `db:"scuteloc" json:"scuteloc"`
	Scutenum           *int      `db:"scutenum" json:"scutenum"`
	Scuteloc2          string    `db:"scuteloc2" json:"scuteloc2"`
	Scutenum2          *int      `db:"scutenum2" json:"scutenum2"`
	Elhv               string    `db:"elhv" json:"elhv"`
	Elcolor            string    `db:"elcolor" json:"elcolor"`
	Erhv               string    `db:"erhv" json:"erhv"`
	Ercolor            string    `db:"ercolor" json:"ercolor"`
	Cwtyn              string    `db:"cwtyn" json:"cwtyn"`
	Dangler            string    `db:"dangler" json:"dangler"`
	Genetic            string    `db:"genetic" json:"genetic"`
	GeneticsVialNumber string    `db:"genetics_vial_number" json:"geneticsVialNumber"`
	Broodstock         *int      `db:"broodstock" json:"broodstock"`
	HatchWild          *int      `db:"hatch_wild" json:"hatchWild"`
	SpeciesId          *int      `db:"species_id" json:"speciesId"`
	Archive            *int      `db:"archive" json:"archive"`
	Head               *int      `db:"head" json:"head"`
	Snouttomouth       *int      `db:"snouttomouth" json:"snouttomouth"`
	Inter              *int      `db:"inter" json:"inter"`
	Mouthwidth         *int      `db:"mouthwidth" json:"mouthwidth"`
	MIb                *int      `db:"m_ib" json:"mIb"`
	LOb                *int      `db:"l_ob" json:"lOb"`
	LIb                *int      `db:"l_ib" json:"lIb"`
	RIb                *int      `db:"r_ib" json:"rIb"`
	ROb                *int      `db:"r_ob" json:"rOb"`
	Anal               *int      `db:"anal" json:"anal"`
	Dorsal             *int      `db:"dorsal" json:"dorsal"`
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
	SeFieldID        string    `db:"se_field_id" json:"seFieldId"`
	Season           string    `db:"season" json:"season"`
	SetDate          string    `db:"setdate" json:"setdate"` // TODO: date not displaying in table
	SetDateTime      time.Time `db:"setdate" json:"setDateTime"`
	Subsample        float64   `db:"subsample" json:"subsample"`
	Subsamplepass    float64   `db:"subsamplepass" json:"subsamplepass"`
	SubsampleROrN    string    `db:"subsample_r_or_n" json:"subsampleROrN"`
	Subsamplen       string    `db:"subsamplen" json:"subsamplen"`
	Recorder         string    `db:"recorder" json:"recorder"`
	Gear             string    `db:"gear" json:"gear"`
	GearType         string    `db:"gear_type" json:"gearType"`
	Temp             float64   `db:"temp" json:"temp"`
	Turbidity        *float64  `db:"turbidity" json:"turbidity"`
	Conductivity     *float64  `db:"conductivity" json:"conductivity"`
	Do               *float64  `db:"do" json:"dissolvedOxygen"`
	Distance         *float64  `db:"distance" json:"distance"`
	Width            *float64  `db:"width" json:"width"`
	Netrivermile     *float64  `db:"netrivermile" json:"netrivermile"`
	Structurenumber  string    `db:"structurenumber" json:"structurenumber"`
	Usgs             string    `db:"usgs" json:"usgs"`
	Riverstage       *float64  `db:"riverstage" json:"riverstage"`
	Discharge        *float64  `db:"discharge" json:"discharge"`
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
	SetSite1         string    `db:"set_site_1" json:"setSite1"`
	SetSite2         string    `db:"set_site_2" json:"setSite2"`
	SetSite3         string    `db:"set_site_3" json:"setSite3"`
	StartTime        string    `db:"starttime" json:"startTime"`
	StartLatitude    float64   `db:"startlatitude" json:"startlatitude"`
	StartLongitude   float64   `db:"startlongitude" json:"startlongitude"`
	StopTime         string    `db:"stoptime" json:"stoptime"`
	StopLatitude     *float64  `db:"stoplatitude" json:"stoplatitude"`
	StopLongitude    *float64  `db:"stop_longitude" json:"stoplongitude"`
	Depth1           *float64  `db:"depth1" json:"depth1"`
	Velocitybot1     *float64  `db:"velocitybot1" json:"velocitybot1"`
	Velocity08_1     *float64  `db:"velocity08_1" json:"velocity081"`
	Velocity02or06_1 *float64  `db:"velocity02or06_1" json:"velocity02or061"`
	Depth2           *float64  `db:"depth2" json:"depth2"`
	Velocitybot2     *float64  `db:"velocitybot2" json:"velocitybot2"`
	Velocity08_2     *float64  `db:"velocity08_2" json:"velocity082"`
	Velocity02or06_2 *float64  `db:"velocity02or06_2" json:"velocity02or062"`
	Depth3           *float64  `db:"depth3" json:"depth3"`
	Velocitybot3     *float64  `db:"velocitybot3" json:"velocitybot3"`
	Velocity08_3     *float64  `db:"velocity08_3" json:"velocity083"`
	Velocity02or06_3 *float64  `db:"velocity02or06_3" json:"velocity02or063"`
	Watervel         *float64  `db:"watervel" json:"watervel"`
	Cobble           *float64  `db:"cobble" json:"cobble"`
	Organic          *float64  `db:"organic" json:"organic"`
	Silt             *float64  `db:"silt" json:"silt"`
	Sand             *float64  `db:"sand" json:"sand"`
	Gravel           *float64  `db:"gravel" json:"gravel"`
	Comments         string    `db:"comments" json:"comments"`
	Complete         *float64  `db:"complete" json:"complete"`
	Checkby          string    `db:"checkby" json:"checkby"`
	NoTurbidity      string    `db:"no_turbidity" json:"noTurbidity"`
	NoVelocity       string    `db:"no_velocity" json:"noVelocity"`
	EditInitials     string    `db:"edit_initials" json:"editInitials"`
	LastEditComment  string    `db:"last_edit_comment" json:"lastEditComment"`
	Project          *int      `db:"PROJECT_ID" json:"project"`
	FieldOffice      string    `db:"FIELDOFFICE" json:"fieldOffice"`
	Segment          *int      `db:"SEGMENT_ID" json:"segment"`
	BkgColor         string    `db:"bkg_color" json:"bkgColor"`
	SuppBkgColor     string    `db:"supp_bkg_color" json:"suppBkgColor"`
	FishCount        int       `db:"fish_count" json:"fishCount"`
	SuppCount        int       `db:"supp_count" json:"suppCount"`
	Bend             int       `db:"bend" json:"bend"`
	BendRn           string    `db:"bendrn" json:"bendrn"`
	BendRiverMile    float64   `db:"bendrivermile" json:"bendrivermile"`
	LastUpdated      time.Time `db:"last_updated" json:"lastUpdated"`
	UploadSessionId  int       `db:"upload_session_id" json:"uploadSessionId"`
	UploadedBy       string    `db:"uploaded_by" json:"uploadedBy"`
	UploadFilename   string    `db:"upload_filename" json:"uploadFilename"`
}

type TelemetryDataEntryWithCount struct {
	Items      []UploadTelemetry `json:"items"`
	TotalCount int               `json:"totalCount"`
}

type UploadTelemetryData struct {
	Items          []UploadTelemetry `json:"items"`
	UploadFilename string            `db:"upload_filename" json:"uploadFilename"`
}

type UploadTelemetry struct {
	TId                int       `db:"t_id" json:"tId"`
	TFid               string    `db:"t_fid" json:"tFid"`
	SeFid              string    `db:"se_fid" json:"seFieldId"`
	SeId               int       `db:"se_id" json:"seId"`
	SiteId             int       `db:"site_id" json:"siteId"`
	Bend               float64   `db:"bend" json:"bend"`
	RadioTagNum        int       `db:"radio_tag_num" json:"radioTagNum"`
	FrequencyIdCode    int       `db:"frequency_id_code" json:"frequencyIdCode"`
	CaptureTime        string    `db:"capture_time" json:"captureDate"`
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
	Checkby            string    `db:"checkby" json:"checkby"`
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

type UnapprovedDataWithCount struct {
	Items      []UnapprovedData `json:"items"`
	TotalCount int              `json:"totalCount"`
}

type UnapprovedData struct {
	Ch                 string `db:"ch" json:"ch"`
	Fp                 string `db:"fp" json:"fp"`
	SegmentDescription string `db:"segment_description" json:"segmentDescription"`
	Bend               int    `db:"bend" json:"bend"`
	MrId               int    `db:"mr_id" json:"mrId"`
	// SetDate            time.Time `db:"setdate" json:"setdate"`
	Subsample      int     `db:"subsampple" json:"subsample"`
	Recorder       string  `db:"recorder" json:"recorder"`
	Checkby        string  `db:"checkby" json:"checkby"`
	NetRiverMile   float64 `db:"netrivermile" json:"netrivermile"`
	SiteId         int     `db:"site_id" json:"siteId"`
	ProjectId      int     `db:"project_id" json:"projectId"`
	SegmentId      int     `db:"segment_id" json:"segmentId"`
	Season         string  `db:"season" json:"season"`
	FieldOffice    string  `db:"fieldoffice" json:"fieldoffice"`
	SampleUnitType string  `db:"sample_unit_type" json:"sampleUnitType"`
	Gear           string  `db:"gear" json:"gear"`
}

type BafiDataWithCount struct {
	Items      []BafiData `json:"items"`
	TotalCount int        `json:"totalCount"`
}

type BafiData struct {
	Psb           string  `db:"psb" json:"psb"`
	SiteId        int     `db:"site_id" json:"siteId"`
	FieldOffice   string  `db:"fieldoffice" json:"fieldoffice"`
	FId           int     `db:"f_id" json:"fId"`
	MrId          int     `db:"mr_id" json:"mrId"`
	MrFid         string  `db:"mr_fid" json:"mrFid"`
	Species       string  `db:"species" json:"species"`
	Recorder      string  `db:"recorder" json:"recorder"`
	Subsample     int     `db:"subsample" json:"subsample"`
	Gear          string  `db:"gear" json:"gear"`
	FishCount     int     `db:"fishcount" json:"fishcount"`
	Year          int     `db:"year" json:"year"`
	SegmentId     int     `db:"segment_id" json:"segmentId"`
	Bend          int     `db:"bend" json:"bend"`
	Bendrn        string  `db:"bendrn" json:"bendrn"`
	BendRiverMile float64 `db:"bendrivermile" json:"bendrivermile"`
	PanelHook     string  `db:"panelhook" json:"panelhook"`
}

type UncheckedDataWithCount struct {
	Items      []UncheckedData `json:"items"`
	TotalCount int             `json:"totalCount"`
}

type UncheckedData struct {
	Cb           string  `db:"cb" json:"cb"`
	Psb          string  `db:"psb" json:"psb"`
	MrID         int     `db:"mr_id" json:"mrID"`
	Subsample    int     `db:"subsample" json:"subsample"`
	Recorder     string  `db:"recorder" json:"recorder"`
	Checkby      string  `db:"checkby" json:"checkby"`
	Netrivermile float64 `db:"netrivermile" json:"netrivermile"`
	SiteID       int     `db:"site_id" json:"siteId"`
	ProjectID    int     `db:"project_id" json:"projectId"`
	SegmentID    int     `db:"segment_id" json:"segmentId"`
	Season       string  `db:"season" json:"season"`
	FieldOffice  string  `db:"fieldoffice" json:"fieldoffice"`
	Gear         string  `db:"gear" json:"gear"`
}

type ErrorCount struct {
	Year  int `db:"year" json:"year"`
	Count int `db:"count(el.el_id)" json:"count"`
}

type OfficeErrorLog struct {
	ElID             int       `db:"el_id" json:"elId"`
	SiteID           int       `db:"site_id" json:"siteId"`
	Year             int       `db:"year" json:"year"`
	ErrorEntryDate   time.Time `db:"error_entry_date" json:"errorEntryDate"`
	WorksheetID      int       `db:"worksheet_id" json:"worksheetId"`
	WorksheetTypeID  int       `db:"worksheet_type_id" json:"worksheetTypeId"`
	FieldID          int       `db:"field_id" json:"fieldId"`
	FormID           int       `db:"form_id" json:"formId"`
	ErrorDescription string    `db:"error_description" json:"errorDescription"`
	ErrorStatus      int       `db:"error_fixed" json:"errorFixed"`
	ErrorUpdateDate  time.Time `db:"error_fixed_date" json:"errorFixedDate"`
}

type DownloadInfo struct {
	Name        string `db:"name" json:"name"`
	DisplayName string `db:"display_name" json:"displayName"`
	LastUpdated string `db:"last_updated" json:"lastUpdated"`
}
