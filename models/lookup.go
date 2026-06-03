package models

// GLOBAL LOOK UP MODELS

type FrequencyId struct {
	FrequencyIdCode int `db:"frequency_id_code" json:"code"`
	FrequencyIdDescription   string `db:"frequency_id_description" json:"description"`
}

type ScuteLocation struct {
	ScuteLocationCode string `db:"scute_location_code" json:"code"`
	ScuteLocationDescription   string `db:"scute_location_description" json:"description"`
}