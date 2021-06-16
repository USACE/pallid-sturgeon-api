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

func (s *PallidSturgeonStore) UploadSiteDatasheetCheck(uploadedBy string, uploadSessionId int) error {
	uploadSiteStmt, err := s.db.Prepare("begin DATA_UPLOAD.uploadSiteDatasheetCheck (:1,:2);  end;")

	//var retVal string
	uploadSiteStmt.Exec(godror.PlSQLArrays, uploadedBy, uploadSessionId)

	//fmt.Println(retVal)

	return err
}

func (s *PallidSturgeonStore) UploadSiteDatasheet(uploadedBy string) error {
	uploadSiteStmt, err := s.db.Prepare("begin DATA_UPLOAD.uploadSiteDatasheet  (:1);  end;")

	uploadSiteStmt.Exec(godror.PlSQLArrays, uploadedBy)

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

func (s *PallidSturgeonStore) UploadFishDatasheetCheck(uploadedBy string, uploadSessionId int) error {
	uploadSiteStmt, err := s.db.Prepare("begin DATA_UPLOAD.uploadFishDatasheetCheck (:1,:2);  end;")

	//var retVal string
	uploadSiteStmt.Exec(godror.PlSQLArrays, uploadedBy, uploadSessionId)

	//fmt.Println(retVal)

	return err
}

func (s *PallidSturgeonStore) UploadFishDatasheet(uploadedBy string) error {
	uploadSiteStmt, err := s.db.Prepare("begin DATA_UPLOAD.uploadFishDatasheet    (:1);  end;")

	uploadSiteStmt.Exec(godror.PlSQLArrays, uploadedBy)

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
		uploadSupplemental.Comments,
		uploadSupplemental.LastUpdated,
		uploadSupplemental.UploadSessionId,
		uploadSupplemental.UploadedBy,
		uploadSupplemental.UploadFilename,
	)

	return err
}

var insertProcedureUploadSql = `insert into upload_procedure (id, f_fid, purpose_code, procedure_date, procedure_start_time, procedure_end_time, procedure_by, 
	antibiotic_injection_ind, photo_dorsal_ind, photo_ventral_ind, photo_left_ind,
	old_radio_tag_num, old_frequency_id, dst_serial_num, dst_start_date, dst_start_time, dst_reimplant_ind, new_radio_tag_num,
	new_frequency_id, sex_code, blood_sample_ind, egg_sample_ind, comments, fish_health_comments,
	eval_location_code, spawn_code, visual_repro_status_code, ultrasound_repro_status_code,
	expected_spawn_year, ultrasound_gonad_length, gonad_condition,
	last_updated, upload_session_id, uploaded_by, upload_filename )                                                        
values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23,:24,:25,:26,:27,:28,:29,:30,:31,:32,:33,:34,:35)`

func (s *PallidSturgeonStore) SaveProcedureUpload(uploadProcedure models.UploadProcedure) error {
	_, err := s.db.Exec(insertProcedureUploadSql,
		uploadProcedure.Id,
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

var insertMrUploadSql = `insert into upload_mr (site_id, site_fid, mr_fid, season, setdate, subsample, subsamplepass, 
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

func (s *PallidSturgeonStore) SaveMrUpload(uploadMr models.UploadMr) error {
	_, err := s.db.Exec(insertMrUploadSql,
		uploadMr.SiteID, uploadMr.SiteFid, uploadMr.MrFid, uploadMr.Season, uploadMr.Setdate,
		uploadMr.Subsample, uploadMr.Subsamplepass, uploadMr.Subsamplen, uploadMr.Recorder,
		uploadMr.Gear, uploadMr.GearType, uploadMr.Temp, uploadMr.Turbidity, uploadMr.Conductivity,
		uploadMr.Do, uploadMr.Distance, uploadMr.Width, uploadMr.Netrivermile, uploadMr.Structurenumber,
		uploadMr.Usgs, uploadMr.Riverstage, uploadMr.Discharge, uploadMr.U1, uploadMr.U2, uploadMr.U3, uploadMr.U4,
		uploadMr.U5, uploadMr.U6, uploadMr.U7, uploadMr.Macro, uploadMr.Meso, uploadMr.Habitatrn, uploadMr.Qc,
		uploadMr.MicroStructure, uploadMr.StructureFlow, uploadMr.StructureMod, uploadMr.SetSite1, uploadMr.SetSite2, uploadMr.SetSite3,
		uploadMr.StartTime, uploadMr.StartLatitude, uploadMr.StartLongitude, uploadMr.StopTime, uploadMr.StopLatitude, uploadMr.StopLongitude,
		uploadMr.Depth1, uploadMr.Velocitybot1, uploadMr.Velocity08_1, uploadMr.Velocity02or06_1,
		uploadMr.Depth2, uploadMr.Velocitybot2, uploadMr.Velocity08_2, uploadMr.Velocity02or06_2,
		uploadMr.Depth3, uploadMr.Velocitybot3, uploadMr.Velocity08_3, uploadMr.Velocity02or06_3,
		uploadMr.Watervel, uploadMr.Cobble, uploadMr.Organic, uploadMr.Silt, uploadMr.Sand, uploadMr.Gravel,
		uploadMr.Comments, uploadMr.LastUpdated, uploadMr.UploadSessionId,
		uploadMr.UploadedBy, uploadMr.UploadFilename, uploadMr.Complete, uploadMr.Checkby,
		uploadMr.NoTurbidity, uploadMr.NoVelocity,
	)

	return err
}

var insertTelemetryFishUploadSql = `insert into upload_telemetry_fish(t_fid, se_fid, bend, radio_tag_num, frequency_id_code, capture_time, capture_latitude, capture_longitude,
	position_confidence, macro_id, meso_id, depth, temp, conductivity, turbidity, silt, sand, gravel, comments, last_updated, upload_session_id, uploaded_by, upload_filename)
	values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23)`

func (s *PallidSturgeonStore) SaveTelemetryFishUpload(uploadTelemetryFish models.UploadTelemetryFish) error {
	_, err := s.db.Exec(insertTelemetryFishUploadSql,
		uploadTelemetryFish.TFid,
		uploadTelemetryFish.SeFid,
		uploadTelemetryFish.Bend,
		uploadTelemetryFish.RadioTagNum,
		uploadTelemetryFish.FrequencyIdCode,
		uploadTelemetryFish.CaptureTime,
		uploadTelemetryFish.CaptureLatitude,
		uploadTelemetryFish.CaptureLongitude,
		uploadTelemetryFish.PositionConfidence,
		uploadTelemetryFish.MacroId,
		uploadTelemetryFish.MesoId,
		uploadTelemetryFish.Depth,
		uploadTelemetryFish.Temp,
		uploadTelemetryFish.Conductivity,
		uploadTelemetryFish.Turbidity,
		uploadTelemetryFish.Silt,
		uploadTelemetryFish.Sand,
		uploadTelemetryFish.Gravel,
		uploadTelemetryFish.Comments,
		uploadTelemetryFish.LastUpdated,
		uploadTelemetryFish.UploadSessionId,
		uploadTelemetryFish.UploadedBy,
		uploadTelemetryFish.UploadFilename,
	)

	return err
}
