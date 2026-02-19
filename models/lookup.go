package models

type BendSelection struct {
	BsId              int    `db:"bs_id" json:"id"`
	BendSelectionCode string `db:"bend_selection_code" json:"code"`
	BendSelectionDesc string `db:"bend_selection_description" json:"description"`
}

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
	MicroStructureCode        string `db:"micro_structure_code" json:"code"`
	MicroStructureDescription string `db:"micro_structure" json:"description"`
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

type SetSite3 struct {
	SsCode        int    `db:"set_site_3_code" json:"code"`
	SsDescription string `db:"set_site_3" json:"description"`
}

type BendRiverMile struct {
	BrmId          int      `db:"brm_id" json:"id"`
	Segment        int      `db:"b_segment" json:"segment"`
	Bend           int      `db:"bend_num" json:"bend"`
	State          string   `db:"state" json:"state"`
	UpperRiverMile *float64 `db:"upper_river_mile" json:"upperRiverMile"`
	LowerRiverMile *float64 `db:"lower_river_mile" json:"lowerRiverMile"`
}