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

type MicroHabitat struct {
	MhId               int    `db:"mh_id" json:"id"`
	MicroStructure     string `db:"micro_structure" json:"microStructure"`
	MicroStructureCode int    `db:"micro_structure_code" json:"microStructureCode"`
	StructureFlow      string `db:"structure_flow" json:"structureFlow"`
	StructureFlowCode  int    `db:"structure_flow_code" json:"structureFlowCode"`
	StructureMod       string `db:"structure_mod" json:"structureMod"`
	StructureModCode   int    `db:"structure_mod_code" json:"structureModCode"`
}

type StructureModLK struct {
	ModDescription string `db:"mod_description" json:"description"`
	ModCode        int    `db:"mod_code" json:"code"`
}

type U7 struct {
	U7Code        string `db:"code" json:"code"`
	U7Description string `db:"description" json:"description"`
}

type LookupData struct {
	BendSelectionData BendSelection `json:"bendSelection"`
	GearCodeData      GearCode      `json:"gearCode"`
	GearTypeData      GearType      `json:"gearType"`
	MacroData         Macro         `json:"macro"`
}
