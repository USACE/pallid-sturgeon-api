package models

// GLOBAL LOOK UP MODELS

type YearLK struct {
	YearId int `db:"year_id" json:"id"`
	Year   int `db:"year" json:"year"`
}

type FieldOfficeLK struct {
	FoId          int    `db:"fo_id" json:"id"`
	FoCode        string `db:"field_office_code" json:"code"`
	FoDescription string `db:"field_office_description" json:"description"`
	State         string `db:"state" json:"state"`
}

type FieldOfficeSegment struct {
	FosId          		int    `db:"fos_id" json:"id"`
	FieldOfficeCode     string `db:"field_office_code" json:"fieldOfficeCode"`
	SegmentCode 		int `db:"segment_code" json:"segmentCode"`
	ProjectCode         int `db:"project_code" json:"projectCode"`
}

type ProjectLK struct {
	ProjectCode        int    `db:"project_code" json:"code"`
	ProjectDescription string `db:"project_description" json:"description"`
}

type SegmentLK struct {
	SegmentId          int    `db:"s_id" json:"id"`
	SegmentCode        int    `db:"segment_code" json:"code"`
	SegmentDescription string `db:"segment_description" json:"description"`
}

type SeasonLK struct {
	SeasonId          int    `db:"s_id" json:"id"`
	SeasonCode        string    `db:"season_code" json:"code"`
	SeasonDescription string `db:"season_description" json:"description"`
}

type SampleUnitTypeLK struct {
	SutCode        string    `db:"sample_unit_type_code" json:"code"`
	SutDescription string `db:"sample_unit_type_description" json:"description"`
}

type BendRiverMile struct {
	BrmId          int      `db:"brm_id" json:"id"`
	Segment        int      `db:"b_segment" json:"segment"`
	Bend           int      `db:"bend_num" json:"bend"`
	BendDescription string `db:"b_desc" json:"bendDescription"`
	State          string   `db:"state" json:"state"`
	UpperRiverMile *float64 `db:"upper_river_mile" json:"upperRiverMile"`
	LowerRiverMile *float64 `db:"lower_river_mile" json:"lowerRiverMile"`
}

type BendSelection struct {
	BsId              int    `db:"bs_id" json:"id"`
	BendSelectionCode string `db:"bend_selection_code" json:"code"`
	BendSelectionDesc string `db:"bend_selection_description" json:"description"`
}

type Chute struct {
	ChuteId          int     `db:"chute_id" json:"id"`
	Segment 	     int 	 `db:"segment_id" json:"segment"`
	ChuteCode 		 string  `db:"chute_code" json:"code"`
	ChuteDescription string  `db:"chute_desc" json:"description"`
	UpperRiverMile 	 float64 `db:"upper_river_mile" json:"upperRiverMile"`
}

type Reach struct {
	ReachId          int     `db:"reach_id" json:"id"`
	Segment 	     int 	 `db:"segment_id" json:"segment"`
	ReachCode 		 string  `db:"reach_code" json:"code"`
	ReachDescription string  `db:"reach_desc" json:"description"`
	UpperRiverMile 	 float64 `db:"upper_river_mile" json:"upperRiverMile"`
}

type Recapture struct {
	Species string `db:"species" json:"species"`
	PitTag string `db:"pit_tag" json:"pitTag"`
}

// MISSOURI RIVER LOOK UP MODELS

type GearCode struct {
	GearId          int    `db:"gear_id" json:"id"`
	Gear            string `db:"gear" json:"gear"`
	GearCode        string `db:"gear_code" json:"code"`
	GearType        string `db:"gear_type" json:"gearType"`
	GearDescription string `db:"gear_description" json:"description"`
	DeploymentType  string `db:"deploymenttype" json:"deploymentType"`
}

type FilteredGearCode struct {
	SgoId           int    `db:"sgo_id" json:"id"`
	FieldOfficeCode string `db:"field_office_code" json:"fieldOfficeCode"`
	SeasonCode      string `db:"season_code" json:"seasonCode"`
	GearCode        string `db:"gear_code" json:"gearCode"`
	Gear            string `db:"gear" json:"gear"`
	ProjectCode     int    `db:"project_code" json:"projectCode"`
}

type GearType struct {
	GtId                int    `db:"gt_id" json:"id"`
	GearTypeCode        string `db:"gear_type_code" json:"code"`
	GearTypeDescription string `db:"gear_type_description" json:"description"`
}

type Macro struct {
	MhId               int    `db:"mh_id" json:"id"`
	HabitatCode        string `db:"habitat_code" json:"code"`
	HabitatDescription string `db:"habitat_description" json:"description"`
}

type MesoLk struct {
	MhId                   int    `db:"mh_id" json:"id"`
	MesoHabitatCode        string `db:"mesohabitat_code" json:"code"`
	MesoHabitatDescription string `db:"mesohabitat_description" json:"description"`
}

type MacroMeso struct {
	MmId             int    `db:"mm_id" json:"id"`
	MacroHabitatCode string `db:"macrohabitat_code" json:"macroHabitatCode"`
	MesoHabitatCode  string `db:"mesohabitat_code" json:"mesoHabitatCode"`
}

type MicroStructure struct {
	MsId                      int    `db:"ms_id" json:"id"`
	MicroStructureCode        string `db:"code" json:"code"`
	MicroStructureDescription string `db:"description" json:"description"`
}

type StructureFlowLK struct {
	SfId                     int    `db:"sf_id" json:"id"`
	StructureFlowCode        string `db:"code" json:"code"`
	StructureFlowDescription string `db:"description" json:"description"`
}

type StructureModLK struct {
	SmId                    int    `db:"sm_id" json:"id"`
	StructureModCode        string `db:"code" json:"code"`
	StructureModDescription string `db:"description" json:"description"`
}

type MicroHabitat struct {
	MhId               int    `db:"mh_id" json:"id"`
	MicroStructure     string `db:"micro_structure" json:"microStructure"`
	MicroStructureCode int    `db:"micro_structure_code" json:"microStructureCode"`
	StructureFlow      string `db:"structure_flow" json:"structureFlow"`
	StructureFlowCode  int    `db:"structure_flow_code" json:"structureFlowCode"`
	StructureMod       string `db:"structure_mod" json:"structureMod"`
	StructureModCode   int    `db:"structure_mod_code" json:"structureModCode"`
}

type U6 struct {
	U6Id          int    `db:"u6_id" json:"id"`
	U6Code        string `db:"code" json:"code"`
	U6Description string `db:"description" json:"description"`
}

type U7 struct {
	U7Code        string `db:"code" json:"code"`
	U7Description string `db:"description" json:"description"`
}

type Estimation struct {
	EstId   int    `db:"coe_id" json:"id"`
	EstCode int    `db:"estimation_code" json:"code"`
	EstDesc string `db:"estimation_description" json:"description"`
}

type MicroSetSite struct {
	MsId               int    `db:"ms_id" json:"id"`
	MicroStructureCode int    `db:"structure_code" json:"microStructureCode"`
	MicroStructureDesc string `db:"micro_structure" json:"microStructureDescription"`
	Ss1Code            int    `db:"set_site_1_code" json:"ss1Code"`
	Ss1Description     string `db:"set_site_1" json:"ss1Description"`
	Ss2Code            int    `db:"set_site_two_code" json:"ss2Code"`
	Ss2Description     string `db:"set_site_two" json:"ss2Description"`
}

type SetSite1LK struct {
	Ss1Id          int    `db:"ss1_id" json:"id"`
	Ss1Code        int    `db:"code" json:"code"`
	Ss1Description string `db:"description" json:"description"`
}

type SetSite2LK struct {
	Ss2Id          int    `db:"ss2_id" json:"id"`
	Ss2Code        int    `db:"code" json:"code"`
	Ss2Description string `db:"description" json:"description"`
}

type SetSite3 struct {
	SsCode        int    `db:"set_site_3_code" json:"code"`
	SsDescription string `db:"set_site_3" json:"description"`
}

type SubsampleType struct {
	StId          int    `db:"st_id" json:"id"`
	StCode        string `db:"code" json:"code"`
	StDescription string `db:"description" json:"description"`
}

// FISH LOOK UP MODELS

type FishCode struct {
	FishId           int     `db:"lk_id" json:"id"`
	CommonName       string  `db:"common_name" json:"commonName"`
	ScientificName   string  `db:"scientific_name" json:"scientificName"`
	AlphaCode        string  `db:"alpha_code" json:"alphaCode"`
	NumericCodes     *int    `db:"numeric_codes" json:"numericCodes"`
	NumericCodesText *string `db:"numeric_codes_txt" json:"numericCodesText"`
}

type FishStructure struct {
	FsCode        string `db:"code" json:"code"`
	FsDescription string `db:"description" json:"description"`
}

type FloyTagPrefix struct {
	FtpId          int    `db:"floy_id" json:"id"`
	FtpCode        string `db:"tag_prefix_code" json:"code"`
	FtpDescription string `db:"tag_prefix_description" json:"description"`
}

type LengthType struct {
	LtId          int    `db:"lk_id" json:"id"`
	LtCode        string `db:"code" json:"code"`
	LtDescription string `db:"description" json:"description"`
}

type MarkRecapture struct {
	MrId          int    `db:"mr_id" json:"id"`
	MrCode        string `db:"mark_recapture_code" json:"code"`
	MrDescription string `db:"mark_recapture_description" json:"description"`
}

type FrequencyId struct {
	FrequencyIdCode 		int 	`db:"frequency_id_code" json:"code"`
	FrequencyIdDescription 	string 	`db:"frequency_id_description" json:"description"`
}

type SpawnBehavior struct {
	SpawnCode			int		`db:"spawn_code" json:"code"`
	SpawnDescription	string 	`db:"spawn_description" json:"description"`
}

type PositionConfidence struct {
	PositionCode 		int 	`db:"position_confidence_code" json:"code"`
	PositionDescription string `db:"position_confidence_description" json:"description"`
}

type SearchType struct {
	SearchTypeCode	string `db:"search_type_code" json:"code"`
	SearchTypeDescription string `db:"search_type_description" json:"description"`
}