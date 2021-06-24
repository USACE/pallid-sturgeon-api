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

var nextUploadSessionIdSql = `SELECT upload_session_seq.nextval from dual`

func (s *PallidSturgeonStore) GetUploadSessionId() (int, error) {
	rows, err := s.db.Query(nextUploadSessionIdSql)

	var nextUploadSessionId int
	for rows.Next() {
		rows.Scan(&nextUploadSessionId)
	}

	return nextUploadSessionId, err
}

var insertUploadSiteSql = `insert into upload_site (site_id, site_fid, site_year, fieldoffice_id, 
	field_office, project_id, project, 
	segment_id, segment, season_id, season, bend, bendrn, bend_river_mile, comments,
	last_updated, upload_session_id, uploaded_by, upload_filename)
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
		uploadFish.FFid,
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
		uploadSearch.SiteID,
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
	sex, stage,  recapture, photo,
	genetic_needs, other_tag_info,
	comments, 
	last_updated, upload_session_id,uploaded_by, upload_filename)
	
	 values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,
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

var insertProcedureUploadSql = `insert into upload_procedure (f_fid, purpose_code, procedure_date, procedure_start_time, procedure_end_time, procedure_by, 
	antibiotic_injection_ind, photo_dorsal_ind, photo_ventral_ind, photo_left_ind,
	old_radio_tag_num, old_frequency_id, dst_serial_num, dst_start_date, dst_start_time, dst_reimplant_ind, new_radio_tag_num,
	new_frequency_id, sex_code, blood_sample_ind, egg_sample_ind, comments, fish_health_comments,
	eval_location_code, spawn_code, visual_repro_status_code, ultrasound_repro_status_code,
	expected_spawn_year, ultrasound_gonad_length, gonad_condition,
	last_updated, upload_session_id, uploaded_by, upload_filename )                                                        
values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23,:24,:25,:26,:27,:28,:29,:30,:31,:32,:33,:34)`

func (s *PallidSturgeonStore) SaveProcedureUpload(uploadProcedure models.UploadProcedure) error {
	_, err := s.db.Exec(insertProcedureUploadSql,
		uploadProcedure.FFid,
		uploadProcedure.PurposeCode,
		uploadProcedure.ProcedurDate,
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
		uploadProcedure.SpawnCode,
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

var insertMoriverUploadSql = `insert into upload_moriver (site_id, site_fid, mr_fid, season, setdate, subsample, subsamplepass, 
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
	uploaded_by, upload_filename, complete, checkby,
	no_turbidity, no_velocity)

	 values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23,:24,
		:25,:26,:27,:28,:29,:30,:31,:32,:33,:34,:35,:36,:37,:38,:39,:40,:41,:42,:43,:44,:45,:46,:47,
		:48,:49,:50,:51,:52,:53,:54,:55,:56,:57,:58,:59,:60,:61,:62,:63,:64,:65,:66,:67,:68,:69,:70,
		:71,:72)`

func (s *PallidSturgeonStore) SaveMoriverUpload(UploadMoriver models.UploadMoriver) error {
	_, err := s.db.Exec(insertMoriverUploadSql,
		UploadMoriver.SiteID, UploadMoriver.SiteFid, UploadMoriver.MrFid, UploadMoriver.Season, UploadMoriver.Setdate,
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
		UploadMoriver.UploadedBy, UploadMoriver.UploadFilename, UploadMoriver.Complete, UploadMoriver.Checkby,
		UploadMoriver.NoTurbidity, UploadMoriver.NoVelocity,
	)

	return err
}

var insertTelemetryUploadSql = `insert into upload_telemetry(t_fid, se_fid, bend, radio_tag_num, frequency_id_code, capture_time, capture_latitude, capture_longitude,
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

	uploadFishStmt, err := s.db.Prepare("begin DATA_UPLOAD.uploadFinal (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12);  end;")
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
