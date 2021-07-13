package models

import "time"

type Season struct {
	ID           int    `db:"id" json:"id"`
	Code         string `db:"code" json:"code"`
	Description  string `db:"description" json:"description"`
	FieldAppFlag string `db:"field_app_flag" json:"fieldAppFlag"`
	ProjectCode  *int   `db:"project_code" json:"projectCode"`
}

type Upload struct {
	SiteUpload         []UploadSite         `json:"siteUpload"`
	FishUpload         []UploadFish         `json:"fishUpload"`
	SearchUpload       []UploadSearch       `json:"searchUpload"`
	ProcedureUpload    []UploadProcedure    `json:"procedureUpload"`
	UploadSupplemental []UploadSupplemental `json:"uploadSupplemental"`
	MoriverUpload      []UploadMoriver      `json:"moriverUpload"`
	TelemetryUpload    []UploadTelemetry    `json:"telemetryUpload"`
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

type UploadFish struct {
	SiteID          int       `db:"site_id" json:"siteId"`
	FFid            string    `db:"f_fid" json:"fFid"`
	MrFid           string    `db:"mr_fid" json:"mrFid"`
	Panelhook       string    `db:"panelhook" json:"panelhook"`
	Bait            string    `db:"bait" json:"bait"`
	Species         string    `db:"species" json:"species"`
	Length          float32   `db:"length" json:"length"`
	Weight          float32   `db:"weight" json:"weight"`
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
	UploadedBy      string    `db:"uploaded_by" json:"uploadedBy"`
	UploadFilename  string    `db:"upload_filename" json:"uploadFilename"`
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
	Temp            float64   `db:"temp" json:"temp"`
	Conductivity    float64   `db:"conductivity" json:"conductivity"`
	LastUpdated     time.Time `db:"last_updated" json:"lastUpdated"`
	UploadSessionId int       `db:"upload_session_id" json:"uploadSessionId"`
	UploadedBy      string    `db:"uploaded_by" json:"uploadedBy"`
	UploadFilename  string    `db:"upload_filename" json:"uploadFilename"`
}

type UploadProcedure struct {
	Id                        int       `db:"id" json:"id"`
	FFid                      string    `db:"f_fid" json:"f_fid"`
	PurposeCode               string    `db:"purpose_code" json:"purposeCode"`
	ProcedurDate              time.Time `db:"procedure_date" json:"procedurDate"`
	ProcedureStartTime        string    `db:"procedure_start_time" json:"procedureStartTime"`
	ProcedureEndTime          string    `db:"procedure_end_time" json:"procedureEndTime"`
	ProcedureBy               string    `db:"procedure_by" json:"procedureBy"`
	AntibioticInjectionInd    int       `db:"antibiotic_injection_ind" json:"antibioticInjectionInd"`
	PhotoDorsalInd            int       `db:"photo_dorsal_ind" json:"photoDorsalInd"`
	PhotoVentralInd           int       `db:"photo_ventral_ind" json:"photoVentralInd"`
	PhotoLeftInd              int       `db:"photo_left_ind" json:"photoLeftInd"`
	OldRadioTagNum            int       `db:"old_radio_tag_num" json:"oldRadioTagNum"`
	OldFrequencyId            int       `db:"old_frequency_id" json:"oldFrequencyId"`
	DstSerialNum              int       `db:"dst_serial_num" json:"dstSerialNum"`
	DstStartDate              time.Time `db:"dst_start_date" json:"dstStartDate"`
	DstStartTime              string    `db:"dst_start_time" json:"dstStartTime"`
	DstReimplantInd           int       `db:"dst_reimplant_ind" json:"dstReimplantInd"`
	NewRadioTagNum            int       `db:"new_radio_tag_num" json:"newRadioTagNum"`
	NewFrequencyId            int       `db:"new_frequency_id" json:"newFrequencyId"`
	SexCode                   string    `db:"sex_code" json:"sexCode"`
	BloodSampleInd            int       `db:"blood_sample_ind" json:"bloodSampleInd"`
	EggSampleInd              int       `db:"egg_sample_ind" json:"eggSampleInd"`
	Comments                  string    `db:"comments" json:"comments"`
	FishHealthComments        string    `db:"fish_health_comments" json:"fishHealthComments"`
	EvalLocationCode          string    `db:"eval_location_code" json:"evalLocationCode"`
	SpawnCode                 string    `db:"spawn_code" json:"spawnCode"`
	VisualReproStatusCode     string    `db:"visual_repro_status_code" json:"visualReproStatusCode"`
	UltrasoundReproStatusCode string    `db:"ultrasound_repro_status_code" json:"ultrasoundReproStatusCode"`
	ExpectedSpawnYear         int       `db:"expected_spawn_year" json:"expectedSpawnYear"`
	UltrasoundGonadLength     float64   `db:"ultrasound_gonad_length" json:"ultrasoundGonadLength"`
	GonadCondition            string    `db:"gonad_condition" json:"gonadCondition"`
	LastUpdated               time.Time `db:"last_updated" json:"lastUpdated"`
	UploadSessionId           int       `db:"upload_session_id" json:"uploadSessionId"`
	UploadedBy                string    `db:"uploaded_by" json:"uploadedBy"`
	UploadFilename            string    `db:"upload_filename" json:"uploadFilename"`
}

type UploadSupplemental struct {
	SiteID             int       `db:"site_id" json:"siteId"`
	FFid               string    `db:"f_fid" json:"fFid"`
	MrFid              string    `db:"mr_fid" json:"mrFid"`
	Tagnumber          string    `db:"tagnumber" json:"tagnumber"`
	Pitrn              string    `db:"pitrn" json:"pitrn"`
	Scuteloc           string    `db:"scuteloc" json:"scuteloc"`
	Scutenum           float64   `db:"scutenum" json:"scutenum"`
	Scuteloc2          string    `db:"scuteloc2" json:"scuteloc2"`
	Scutenum2          float64   `db:"scutenum2" json:"scutenum2"`
	Elhv               string    `db:"elhv" json:"elhv"`
	Elcolor            string    `db:"elcolor" json:"elcolor"`
	Erhv               string    `db:"erhv" json:"erhv"`
	Ercolor            string    `db:"ercolor" json:"ercolor"`
	Cwtyn              string    `db:"cwtyn" json:"cwtyn"`
	Dangler            string    `db:"dangler" json:"dangler"`
	Genetic            string    `db:"genetic" json:"genetic"`
	GeneticsVialNumber string    `db:"genetics_vial_number" json:"geneticsVialNumber"`
	Broodstock         float64   `db:"broodstock" json:"broodstock"`
	HatchWild          float64   `db:"hatch_wild" json:"hatchWild"`
	SpeciesId          int       `db:"species_id" json:"speciesId"`
	Archive            int       `db:"archive" json:"archive"`
	Head               float64   `db:"head" json:"head"`
	Snouttomouth       float64   `db:"snouttomouth" json:"snouttomouth"`
	Inter              float64   `db:"inter" json:"inter"`
	Mouthwidth         float64   `db:"mouthwidth" json:"mouthwidth"`
	MIb                float64   `db:"m_ib" json:"mIb"`
	LOb                float64   `db:"l_ob" json:"lOb"`
	LIb                float64   `db:"l_ib" json:"lIb"`
	RIb                float64   `db:"r_ib" json:"rIb"`
	ROb                float64   `db:"r_ob" json:"rOb"`
	Anal               float64   `db:"anal" json:"anal"`
	Dorsal             float64   `db:"dorsal" json:"dorsal"`
	Status             string    `db:"status" json:"status"`
	HatcheryOrigin     string    `db:"hatchery_origin" json:"hatcheryOrigin"`
	Sex                string    `db:"sex" json:"sex"`
	Stage              string    `db:"stage" json:"stage"`
	Recapture          string    `db:"recapture" json:"recapture"`
	Photo              string    `db:"photo" json:"photo"`
	GeneticNeeds       string    `db:"genetic_needs" json:"geneticNeeds"`
	OtherTagInfo       string    `db:"other_tag_info" json:"otherTagInfo"`
	Comments           string    `db:"comments" json:"comments"`
	LastUpdated        time.Time `db:"last_updated" json:"lastUpdated"`
	UploadSessionId    int       `db:"upload_session_id" json:"uploadSessionId"`
	UploadedBy         string    `db:"uploaded_by" json:"uploadedBy"`
	UploadFilename     string    `db:"upload_filename" json:"uploadFilename"`
}

type UploadMoriver struct {
	SiteID           int       `db:"site_id" json:"siteId"`
	SiteFid          string    `db:"site_fid" json:"siteFid"`
	MrFid            string    `db:"mr_fid" json:"mrFid"`
	Season           string    `db:"season" json:"season"`
	Setdate          time.Time `db:"setdate" json:"setdate"`
	Subsample        float64   `db:"subsample" json:"subsample"`
	Subsamplepass    float64   `db:"subsamplepass" json:"subsamplepass"`
	Subsamplen       string    `db:"subsamplen" json:"subsamplen"`
	Recorder         string    `db:"recorder" json:"recorder"`
	Gear             string    `db:"gear" json:"gear"`
	GearType         string    `db:"gear_type" json:"gearType"`
	Temp             float64   `db:"temp" json:"temp"`
	Turbidity        float64   `db:"turbidity" json:"turbidity"`
	Conductivity     float64   `db:"conductivity" json:"conductivity"`
	Do               float64   `db:"do" json:"do"`
	Distance         float64   `db:"distance" json:"distance"`
	Width            float64   `db:"width" json:"width"`
	Netrivermile     float64   `db:"netrivermile" json:"netrivermile"`
	Structurenumber  string    `db:"structurenumber" json:"structurenumber"`
	Usgs             string    `db:"usgs" json:"usgs"`
	Riverstage       float64   `db:"riverstage" json:"riverstage"`
	Discharge        float64   `db:"discharge" json:"discharge"`
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
	StartLatitude    float64   `db:"startlatitude" json:"startLatitude"`
	StartLongitude   float64   `db:"startlongitude" json:"startLongitude"`
	StopTime         string    `db:"stoptime" json:"stopTime"`
	StopLatitude     float64   `db:"stoplatitude" json:"stopLatitude"`
	StopLongitude    float64   `db:"stop_longitude" json:"stopLongitude"`
	Depth1           float64   `db:"depth1" json:"depth1"`
	Velocitybot1     float64   `db:"velocitybot1" json:"velocitybot1"`
	Velocity08_1     float64   `db:"velocity08_1" json:"velocity08_1"`
	Velocity02or06_1 float64   `db:"velocity02or06_1" json:"velocity02or06_1"`
	Depth2           float64   `db:"depth2" json:"depth2"`
	Velocitybot2     float64   `db:"velocitybot2" json:"velocitybot2"`
	Velocity08_2     float64   `db:"velocity08_2" json:"velocity08_2"`
	Velocity02or06_2 float64   `db:"velocity02or06_2" json:"velocity02or06_2"`
	Depth3           float64   `db:"depth3" json:"depth3"`
	Velocitybot3     float64   `db:"velocitybot3" json:"velocitybot3"`
	Velocity08_3     float64   `db:"velocity08_3" json:"velocity08_3"`
	Velocity02or06_3 float64   `db:"velocity02or06_3" json:"velocity02or06_3"`
	Watervel         float64   `db:"watervel" json:"watervel"`
	Cobble           float64   `db:"cobble" json:"cobble"`
	Organic          float64   `db:"organic" json:"organic"`
	Silt             float64   `db:"silt" json:"silt"`
	Sand             float64   `db:"sand" json:"sand"`
	Gravel           float64   `db:"gravel" json:"gravel"`
	Comments         string    `db:"comments" json:"comments"`
	Complete         float64   `db:"complete" json:"complete"`
	Checkby          string    `db:"checkby" json:"checkby"`
	NoTurbidity      string    `db:"no_turbidity" json:"noTurbidity"`
	NoVelocity       string    `db:"no_velocity" json:"noVelocity"`
	LastUpdated      time.Time `db:"last_updated" json:"lastUpdated"`
	UploadSessionId  int       `db:"upload_session_id" json:"uploadSessionId"`
	UploadedBy       string    `db:"uploaded_by" json:"uploadedBy"`
	UploadFilename   string    `db:"upload_filename" json:"uploadFilename"`
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
