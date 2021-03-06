create or replace PACKAGE BODY          data_upload IS
/*------------------------------------------------------------------------------
Description: This package contains procedures for file upload

Name               Date       Modification
------------------ ---------- ------------
------------------------------------------------------------------------------*/


PROCEDURE uploadDSSearchDatasheet(p_user IN upload_supplemental.uploaded_by%TYPE,
                                                                                p_cnt OUT number,
                                                                                p_sfidMatch OUT number,
                                                                                p_fileBrowseSup IN varchar2 default null
                                                                                ) IS
BEGIN
 NULL;
END;


PROCEDURE uploadMRdatasheet (
                                p_user IN upload_mr.uploaded_by%TYPE,
                                p_complete IN upload_mr.complete%TYPE,
                                p_checkby IN upload_mr.checkby%TYPE,
                                p_cnt OUT number,
                                p_mrfidMatch OUT number,
                                p_fileBrowseMR IN varchar2 default null
                            ) IS

v_blob_data       BLOB;
v_blob_len        NUMBER;
v_position        NUMBER;
v_raw_chunk       RAW(10000);
v_char      CHAR(1);
c_chunk_len   number := 1;
v_line        VARCHAR2 (32767) := NULL;
v_data_array      wwv_flow_global.vc_arr2;

v_filename varchar2(200);

-- added for ignoring first row
v_line_count number;

l_session number default 0;
l_max_mr_id number default 0;
-- l_site_id number default 0;
l_mr_max_session number default 0;
l_site_max_session number default 0;

l_mrDuplicate number default 0;
l_mrDuplicateTot number default 0;


BEGIN

-- set session number from sequence
l_session := mr_load_session_seq.nextval;

p_cnt := 0;
p_mrfidMatch := 0;

-- test

-- Read data from wwv_flow_files
 IF p_user is not null THEN

--    select blob_content into v_blob_data
--    from wwv_flow_files 
--    where last_updated = (select max(last_updated) from wwv_flow_files where UPDATED_BY = p_user)
--         and id = (select max(id) from wwv_flow_files where updated_by = p_user);
--         
--    select filename into v_filename
--    from wwv_flow_files 
--    where last_updated = (select max(last_updated) from wwv_flow_files where UPDATED_BY = p_user)
--         and id = (select max(id) from wwv_flow_files where updated_by = p_user);
         
    select blob_content into v_blob_data
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseMR;
         
    select filename into v_filename
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseMR;

 END IF;

v_blob_len := dbms_lob.getlength(v_blob_data);
v_position := 1;

v_line_count := 0;

 
-- Read and convert binary to char
WHILE ( v_position <= v_blob_len ) LOOP
 v_raw_chunk := dbms_lob.substr(v_blob_data,c_chunk_len,v_position);
 v_char :=  chr(hex_to_decimal(rawtohex(v_raw_chunk)));
 v_line := v_line || v_char;
 v_position := v_position + c_chunk_len;

-- When a whole line is retrieved
     IF v_char = CHR(10) THEN

    -- added for ignoring first row   
       v_line_count := v_line_count + 1;

    -- Convert comma to : to use wwv_flow_utilities
       v_line := REPLACE (v_line, ':', '_');
       v_line := REPLACE (v_line, ',', ':');

    -- Convert each column separated by : into array of data
       v_data_array := wwv_flow_utilities.string_to_table (v_line);



       IF v_line_count > 1 THEN

    -- Insert data into target table
       EXECUTE IMMEDIATE 'insert into upload_mr (site_id, site_fid, mr_fid, season, setdate, subsample, subsamplepass, 
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
            :71,:72)'
         USING 
         to_number(v_data_array(1)),
         v_data_array(2),
         v_data_array(3),
         v_data_array(4),
         to_date(v_data_array(5),'mm/dd/yyyy'),
         to_number(v_data_array(6)),
         to_number(v_data_array(7)),
         v_data_array(8),
         v_data_array(9),
         
         v_data_array(10),
         v_data_array(11),
         to_number(v_data_array(12)),
         to_number(v_data_array(13)),
         to_number(v_data_array(14)),
         to_number(v_data_array(15)),
         
         to_number(v_data_array(16)),
         to_number(v_data_array(17)),
         
         to_number(v_data_array(18)),
         v_data_array(19),
         
         v_data_array(20),
         to_number(v_data_array(21)),
         to_number(v_data_array(22)),
         
         v_data_array(23),
         v_data_array(24),
         v_data_array(25),
         v_data_array(26),
         v_data_array(27),
         v_data_array(28),
         v_data_array(29),
         
         v_data_array(30),
         v_data_array(31),
         v_data_array(32),
         v_data_array(33),
         
         v_data_array(34),
         v_data_array(35),
         v_data_array(36),
         v_data_array(37),
         v_data_array(38),
         v_data_array(39),
         
         v_data_array(40),
         to_number(v_data_array(41)),
         to_number(v_data_array(42)),
         v_data_array(43),
         to_number(v_data_array(44)),
         to_number(v_data_array(45)),
         
         to_number(v_data_array(46)),
         to_number(v_data_array(47)),
         to_number(v_data_array(48)),
         to_number(v_data_array(49)),
         
         to_number(v_data_array(50)),
         to_number(v_data_array(51)),
         to_number(v_data_array(52)),
         to_number(v_data_array(53)),
         
         to_number(v_data_array(54)),
         to_number(v_data_array(55)),
         to_number(v_data_array(56)),
         to_number(v_data_array(57)),
         
         to_number(v_data_array(58)),
         to_number(v_data_array(59)),
         to_number(v_data_array(60)),
         to_number(v_data_array(61)),
         to_number(v_data_array(62)),
         to_number(v_data_array(63)),         
         
         v_data_array(66),  --  v_data_array(66),
         sysdate,
         l_session,
         p_user,
         v_filename,
         p_complete,
         p_checkby,
         v_data_array(64),
         v_data_array(65);
         
         commit;

       END IF; -- v_line_count > 1
         
        -- Clear out
           v_line := NULL;
           
     END IF;
  
  END LOOP;

  v_line_count := 0; 
  
  --  get latest mr_id
  SELECT distinct max(mr_id)
    into l_max_mr_id
  FROM upload_mr;
  
  -- get site_id loaded for this upload
--  SELECT SITE_ID
--    INTO l_site_id
--  FROM upload_mr
--  where mr_id = l_max_mr_id;
  
    select distinct max(upload_session_id)
      into l_mr_max_session
    from upload_mr;

    select distinct max(upload_session_id)
        into l_site_max_session
    from upload_sites;
            
            
    -- if uploading a MR data sheet associated with new site
    FOR x IN
        (SELECT SITE_ID
         FROM upload_MR
         WHERE upload_session_id = l_mr_max_session)
     LOOP
        
        IF x.SITE_ID IS NULL OR x.SITE_ID = 0 THEN
        
             update upload_mr
             set site_id = 
                 (SELECT distinct max(ds.site_id) 
                 FROM upload_sites us, ds_sites ds
                 WHERE US.SITE_FID = DS.SITE_FID
                  AND upload_mr.site_fid = US.SITE_FID
                  AND US.UPLOAD_SESSION_ID = l_site_max_session
                  )
             WHERE upload_session_id = l_mr_max_session
                AND site_id = x.site_id;
            
         
        END IF;
     
     END LOOP;
     
  -- Loop through upload_mr records for this session
  FOR Y in (
    SELECT mr_fid
    FROM upload_mr
    WHERE uploaded_by = p_user
     and upload_session_id = l_session)
  
  LOOP
  
    SELECT count(*)
        INTO l_mrDuplicate
    FROM ds_moriver
    WHERE mr_fid = y.mr_fid;
    
    l_mrDuplicateTot := l_mrDuplicateTot +l_mrDuplicate; -- total mr_fid that match 
  
  END LOOP;
  
  -- If there are mr_fid that already exist in DS_MORIVER then don't upload any records to DS_MORIVER
  IF l_mrDuplicateTot = 0 THEN
  
      -- insert into ds_moriver
      INSERT into ds_moriver (site_id, mr_fid, season, setdate, subsample, subsamplepass, comments, last_updated,
            uploaded_by, subsamplen, recorder, gear, gear_type, temp, turbidity, conductivity, do,
            distance, width, netrivermile, structurenumber, usgs, riverstage, discharge,
            u1, u2, u3, u4, u5, u6, u7, macro, meso, habitatrn, qc,
            micro_structure, structure_flow, structure_mod, set_site_1, set_site_2, set_site_3,
            micro,
            starttime, startlatitude, startlongitude, stoptime, stoplatitude, stoplongitude,
            depth1, velocitybot1, velocity08_1, velocity02or06_1,
            depth2, velocitybot2, velocity08_2, velocity02or06_2,
            depth3, velocitybot3, velocity08_3, velocity02or06_3,
            watervel, cobble, organic, silt, sand, gravel, upload_filename,
            complete, checkby, 
            no_turbidity, no_velocity, 
            upload_session_id
            )
        SELECT site_id, mr_fid, season, setdate, subsample, subsamplepass, comments, last_updated, uploaded_by,
            subsamplen, recorder, gear, gear_type, temp, turbidity, conductivity, do,
            distance, width, netrivermile, structurenumber, usgs, riverstage, discharge,
            u1, u2, u3, u4, u5, u6, u7, macro, meso, habitatrn, qc,
            micro_structure, structure_flow, structure_mod, set_site_1, set_site_2, set_site_3,
            micro_structure || structure_flow || structure_mod || set_site_1 || set_site_2 || set_site_3, 
            replace(starttime,'_',':') as start_time, startlatitude, startlongitude, 
            replace(stoptime,'_',':') as stop_time, stoplatitude, stoplongitude,
            depth1, velocitybot1, velocity08_1, velocity02or06_1,
            depth2, velocitybot2, velocity08_2, velocity02or06_2,
            depth3, velocitybot3, velocity08_3, velocity02or06_3,
            watervel, cobble, organic, silt, sand, gravel, upload_filename,
            complete, checkby, 
            no_turbidity, no_velocity, 
            upload_session_id
        FROM UPLOAD_MR
        WHERE uploaded_by = p_user
         and upload_session_id = l_session;
         
      COMMIT;
  
  ELSE
  
    -- number of mr_fid that already exist in DS_MORIVER
    p_mrfidMatch := l_mrDuplicateTot;
    
  END IF;
  
  -- count how many new mr records were added to DS_MORIVER
  select count(*)
   into p_cnt
  from ds_moriver
  where upload_session_id = l_session;
  
  
   -- clear out variables
    l_site_max_session := 0;
    l_mr_max_session := 0;
    
    l_mrDuplicate := 0;
    l_mrDuplicateTot := 0;
 
 
-- Delete file from wwv_flows
-- DELETE FROM wwv_flow_files
-- WHERE id = (select max(id) from wwv_flow_files where updated_by = p_user);
 
 DELETE FROM APEX_APPLICATION_FILES 
 where name = p_fileBrowseMR;


 l_session := 0;
 
END uploadMRdatasheet;


PROCEDURE uploadFishDatasheet (
                                p_user IN upload_fish.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_ffidMatch OUT number,
                                p_fileBrowseF IN varchar2 default null
                            ) IS

v_blob_data       BLOB;
v_blob_len        NUMBER;
v_position        NUMBER;
v_raw_chunk       RAW(10000);
v_char      CHAR(1);
c_chunk_len   number := 1;
v_line        VARCHAR2 (32767) := NULL;
v_data_array      wwv_flow_global.vc_arr2;

v_filename varchar2(200);

-- added for ignoring first row
v_line_count number;

l_session number default 0;
 l_fish_max_session number default 0;
 l_mr_max_session number default 0;
 
l_fDuplicate number default 0;
l_fDuplicateTot number default 0;


BEGIN

-- set session number from sequence
-- l_session := mr_load_session_seq.nextval;
l_session := fish_load_session_seq.nextval;

p_cnt := 0;
p_ffidMatch := 0;

-- Read data from wwv_flow_files
 IF p_user is not null THEN
 
--    select blob_content into v_blob_data
--    from wwv_flow_files 
--    where last_updated = (select max(last_updated) from wwv_flow_files where UPDATED_BY = p_user)
--         and id = (select max(id) from wwv_flow_files where updated_by = p_user);
--         
--    select filename into v_filename
--    from wwv_flow_files 
--    where last_updated = (select max(last_updated) from wwv_flow_files where UPDATED_BY = p_user)
--         and id = (select max(id) from wwv_flow_files where updated_by = p_user);
         
    select blob_content into v_blob_data
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseF;
         
    select filename into v_filename
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseF;
         
         
    
     
 END IF;  
     
     
 
v_blob_len := dbms_lob.getlength(v_blob_data);
v_position := 1;

v_line_count := 0;

 
-- Read and convert binary to char
WHILE ( v_position <= v_blob_len ) LOOP
 v_raw_chunk := dbms_lob.substr(v_blob_data,c_chunk_len,v_position);
 v_char :=  chr(hex_to_decimal(rawtohex(v_raw_chunk)));
 v_line := v_line || v_char;
 v_position := v_position + c_chunk_len;

-- When a whole line is retrieved
     IF v_char = CHR(10) THEN

    -- added for ignoring first row   
       v_line_count := v_line_count + 1;

    -- Convert comma to : to use wwv_flow_utilities
       v_line := REPLACE (v_line, ',', ':');

    -- Convert each column separated by : into array of data
       v_data_array := wwv_flow_utilities.string_to_table (v_line);



       IF v_line_count > 1 THEN

    -- Insert data into target table
       EXECUTE IMMEDIATE 'insert into upload_fish (site_id, f_fid, mr_fid, panelhook, bait, species, length, weight,
        fishcount, fin_curl, otolith, rayspine, scale, ftprefix, ftnum, ftmr,
        comments, last_updated, upload_session_id,uploaded_by, upload_filename)
         values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21)'
         USING 
         to_number(v_data_array(1)),
         v_data_array(2),
         v_data_array(3),
         v_data_array(4),
         v_data_array(5),
         v_data_array(6),
         to_number(v_data_array(7)),
         to_number(v_data_array(8)),         
         to_number(v_data_array(9)),
         v_data_array(10), -- fin_curl - 3/24/16
         v_data_array(11),
         v_data_array(12),
         v_data_array(13),
         v_data_array(14),
         v_data_array(15),
         v_data_array(16),         
         v_data_array(17),
         sysdate,
         l_session,
         p_user,
         v_filename;
         
         commit;

       END IF; -- v_line_count > 1
         
        -- Clear out
           v_line := NULL;
           
     END IF;
  
  END LOOP;

  v_line_count := 0; 
  

select distinct max(upload_session_id)
    into l_mr_max_session
from upload_mr;

select distinct max(upload_session_id)
    into l_fish_max_session
from upload_fish;

 update upload_fish
 set (mr_id, complete, checkby) = 
     (SELECT mr_id, complete, checkby 
     FROM upload_mr
     WHERE upload_fish.MR_FID = upload_mr.MR_FID
      and upload_mr.UPLOAD_SESSION_ID = l_mr_max_session )   
 WHERE upload_session_id = l_fish_max_session;



 -- Loop through upload_fish records for this session
  FOR Y in (
    SELECT f_fid
    FROM upload_fish
    WHERE uploaded_by = p_user
     and upload_session_id = l_session)
  
  LOOP
  
    SELECT count(*)
        INTO l_fDuplicate
    FROM ds_fish
    WHERE f_fid = y.f_fid;
    
    l_fDuplicateTot := l_fDuplicateTot +l_fDuplicate; -- total f_fid that match 
  
  END LOOP;
  
  -- If there are f_fid that already exist in DS_FISH then don't upload any records to DS_FISH
  IF l_fDuplicateTot = 0 THEN
   
      -- insert into ds_fish
      INSERT into ds_fish (mr_id, f_fid, panelhook, bait, species, length, weight,
        last_updated, uploaded_by,fishcount, otolith, rayspine, scale, ftprefix, ftnum, ftmr,
        upload_filename, complete, checkby, upload_session_id)
        SELECT mr_id, f_fid, panelhook, bait, species, length, weight, last_updated, uploaded_by,
            fishcount, otolith, rayspine, scale, ftprefix, ftnum, ftmr, upload_filename,
            complete, checkby, upload_session_id
        FROM UPLOAD_FISH
        WHERE uploaded_by = p_user
         and upload_session_id = l_session
         and species IS NOT NULL;   -- added to remove entries where no species was entered 
         
         
      COMMIT;
  
  ELSE
  
    -- number of mr_fid that already exist in DS_MORIVER
    p_ffidMatch := l_fDuplicateTot;
    
  END IF;
  
  -- count how many new mr records were added to DS_MORIVER
  select count(*)
   into p_cnt
  from ds_fish
  where upload_session_id = l_session; 
  
 
-- clear out variables
l_fish_max_session := 0;
l_mr_max_session := 0;

l_fDuplicate := 0;
l_fDuplicateTot := 0;
 
-- Delete file from wwv_flows
-- DELETE FROM wwv_flow_files
-- WHERE id = (select max(id) from wwv_flow_files where updated_by = p_user);
 
 DELETE from APEX_APPLICATION_FILES 
 where name = p_fileBrowseF;
 
 


 l_session := 0;
 
END uploadFishDatasheet;

---------------------
PROCEDURE uploadSuppDatasheet (
                                p_user IN upload_supplemental.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_sfidMatch OUT number,
                                p_fileBrowseSup IN varchar2 default null
                            ) IS

v_blob_data       BLOB;
v_blob_len        NUMBER;
v_position        NUMBER;
v_raw_chunk       RAW(10000);
v_char      CHAR(1);
c_chunk_len   number       := 1;
v_line        VARCHAR2 (32767)        := NULL;
v_data_array      wwv_flow_global.vc_arr2;

v_filename varchar2(200);

-- added for ignoring first row
v_line_count number;

l_session number default 0;
l_fish_max_session number default 0;
l_mr_max_session number default 0;
l_supp_max_session number default 0;

l_sDuplicate number default 0;
l_sDuplicateTot number default 0;

BEGIN

-- set session number from sequence
l_session := supp_load_session_seq.nextval;

p_cnt := 0;
p_sfidMatch := 0;

-- Read data from wwv_flow_files
 IF p_user is not null THEN

--    select blob_content into v_blob_data
--    from wwv_flow_files 
--    where last_updated = (select max(last_updated) from wwv_flow_files where UPDATED_BY = p_user)
--         and id = (select max(id) from wwv_flow_files where updated_by = p_user);
--         
--    select filename into v_filename
--        from wwv_flow_files 
--        where last_updated = (select max(last_updated) from wwv_flow_files where UPDATED_BY = p_user)
--             and id = (select max(id) from wwv_flow_files where updated_by = p_user);  
             
    select blob_content into v_blob_data
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseSup;
         
    select filename into v_filename
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseSup;
             
 END IF;  

 
v_blob_len := dbms_lob.getlength(v_blob_data);
v_position := 1;

v_line_count := 0;

 
-- Read and convert binary to char
WHILE ( v_position <= v_blob_len ) LOOP
 v_raw_chunk := dbms_lob.substr(v_blob_data,c_chunk_len,v_position);
 v_char :=  chr(hex_to_decimal(rawtohex(v_raw_chunk)));
 v_line := v_line || v_char;
 v_position := v_position + c_chunk_len;

-- When a whole line is retrieved
 IF v_char = CHR(10) THEN

-- added for ignoring first row   
   v_line_count := v_line_count + 1;

-- Convert comma to : to use wwv_flow_utilities
   v_line := REPLACE (v_line, ',', ':');

-- Convert each column separated by : into array of data
   v_data_array := wwv_flow_utilities.string_to_table (v_line);



  IF v_line_count > 1 THEN

-- Insert data into target table
   EXECUTE IMMEDIATE 'insert into upload_supplemental (site_id, f_fid, mr_fid, 
        tagnumber, pitrn, cwtyn, dangler,
        scuteloc, scutenum, elhv, elcolor, erhv, ercolor, genetic, genetics_vial_number,
        head, snouttomouth, inter, mouthwidth, m_ib,
        l_ob, l_ib, r_ib, r_ob, anal, dorsal,
        sex, stage, status, hatchery_origin, recapture, photo,
        genetic_needs, other_tag_info,
        comments, last_updated, upload_session_id,uploaded_by, upload_filename)
         values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,
            :21,:22,:23,:24,:25,:26,:27,:28,:29,:30,:31,:32,:33,:34,:35,:36,:37,:38,:39)'
         USING 
         to_number(v_data_array(1)),
         -- v_data_array(2),
         v_data_array(2),
         v_data_array(3),
         v_data_array(4),         
         v_data_array(5),
         v_data_array(6),
         v_data_array(7),
         
         v_data_array(8),
         to_number(v_data_array(9)),
         
         v_data_array(10),
         v_data_array(11),
         v_data_array(12),
         v_data_array(13),
         v_data_array(14),
         v_data_array(15),
         
         to_number(v_data_array(16)),
         to_number(v_data_array(17)),
         to_number(v_data_array(18)),
         to_number(v_data_array(19)),
         to_number(v_data_array(20)),         
         to_number(v_data_array(21)),
         to_number(v_data_array(22)),
         to_number(v_data_array(23)),
         to_number(v_data_array(24)),
         to_number(v_data_array(25)),
         to_number(v_data_array(26)),
         
         v_data_array(27),
         v_data_array(28),
         v_data_array(29),
         v_data_array(30),
         v_data_array(31),
         v_data_array(32),
         v_data_array(33),
         v_data_array(34),
         
         v_data_array(35),         
         sysdate,
         l_session,
         p_user,
         v_filename;

  END IF; -- v_line_count > 1
 
-- Clear out
   v_line := NULL;
   
 END IF;
  
 END LOOP;

  v_line_count := 0;

-- assign corresponding mr_datashee mr_id to supp records mr_id

select distinct max(upload_session_id)
    into l_supp_max_session
from upload_supplemental;

select distinct max(upload_session_id)
    into l_mr_max_session
from upload_mr;

select distinct max(upload_session_id)
    into l_fish_max_session
from upload_fish;

 update upload_supplemental
 set (mr_id,complete,checkby) = 
     (SELECT mr_id, complete, checkby
     FROM upload_mr
     WHERE upload_supplemental.MR_FID = upload_mr.MR_FID
      and upload_mr.UPLOAD_SESSION_ID = l_mr_max_session ),
      
     f_id = 
     (SELECT f_id 
     FROM upload_fish
     WHERE upload_supplemental.F_FID = upload_fish.F_FID
      and upload_fish.UPLOAD_SESSION_ID = l_fish_max_session )
      
 WHERE upload_session_id = l_supp_max_session;
 
 
 -- Loop through upload_fish records for this session
  FOR Y in (
    SELECT f_fid
    FROM upload_supplemental
    WHERE uploaded_by = p_user
     and upload_session_id = l_session)
  
  LOOP
  
    SELECT count(*)
        INTO l_sDuplicate
    FROM ds_supplemental
    WHERE f_fid = y.f_fid;
    
    l_sDuplicateTot := l_sDuplicateTot + l_sDuplicate; -- total f_fid that match 
  
  END LOOP;
 
 
  -- If there are f_fid that already exist in DS_FISH then don't upload any records to DS_FISH
  IF l_sDuplicateTot = 0 THEN

      -- insert into ds_supplemental
      INSERT into ds_supplemental (mr_id, s_id, f_fid, tagnumber, pitrn, cwtyn, dangler,
        last_updated, uploaded_by, scuteloc, scutenum,
        elhv, elcolor, erhv, ercolor, genetic, genetics_vial_number,
        head, snouttomouth, inter, mouthwidth, m_ib, l_ob, l_ib, r_ib, r_ob, anal, dorsal,
        sex, stage, status, hatchery_origin, recapture, photo,
        genetic_needs, other_tag_info,comments, upload_filename, complete, checkby, upload_session_id)
        SELECT mr_id, s_id, f_fid, tagnumber, pitrn, cwtyn, dangler, last_updated, uploaded_by,
            scuteloc, scutenum, elhv, elcolor, erhv, ercolor, genetic, genetics_vial_number,
            head, snouttomouth, inter, mouthwidth, m_ib, l_ob, l_ib, r_ib, r_ob, anal, dorsal,
            sex, stage, status, hatchery_origin, recapture, photo,
            genetic_needs, other_tag_info, comments, upload_filename, complete, checkby,upload_session_id
        FROM UPLOAD_SUPPLEMENTAL
        WHERE uploaded_by = p_user
         and upload_session_id = l_session;
         
      COMMIT;
      
  ELSE
  
    -- number of mr_fid that already exist in DS_MORIVER
    p_sfidMatch := l_sDuplicateTot;
    
  END IF;
  
  -- count how many new mr records were added to DS_MORIVER
  select count(*)
   into p_cnt
  from ds_supplemental
  where upload_session_id = l_session; 

-- clear out variables
l_fish_max_session := 0;
l_mr_max_session := 0;
l_supp_max_session := 0;

l_sDuplicate := 0;
l_sDuplicateTot := 0;
 
-- Delete file from wwv_flows
-- DELETE FROM wwv_flow_files
-- WHERE id = (select max(id) from wwv_flow_files where updated_by = p_user);
 
 DELETE from APEX_APPLICATION_FILES 
 where name = p_fileBrowseSup; 
 
 l_session := 0;
 
 
END uploadSuppDatasheet;


PROCEDURE uploadSiteDatasheet (
                                p_user IN upload_sites.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_siteMatch OUT number,
                                p_fileBrowseSite IN varchar2 default null
                            ) IS

v_blob_data       BLOB;
v_blob_len        NUMBER;
v_position        NUMBER;
v_raw_chunk       RAW(10000);
v_char      CHAR(1);
c_chunk_len   number       := 1;
v_line        VARCHAR2 (32767)        := NULL;
v_data_array      wwv_flow_global.vc_arr2;

v_filename varchar2(200);

-- added for ignoring first row
v_line_count number;

l_session number default 0;
l_fish_max_session number default 0;
l_mr_max_session number default 0;
l_supp_max_session number default 0;

l_siteMatch number default 0;
l_siteMatch_tot number default 0;


BEGIN

-- set session number from sequence
l_session := site_load_session_seq.nextval;

p_siteMatch := l_siteMatch_tot;

-- Read data from wwv_flow_files
IF p_user is not null THEN

--    select blob_content into v_blob_data
--    from wwv_flow_files 
--    where last_updated = (select max(last_updated) from wwv_flow_files where UPDATED_BY = p_user)
--         and id = (select max(id) from wwv_flow_files where updated_by = p_user);
--         
--    select filename into v_filename
--    from wwv_flow_files 
--    where last_updated = (select max(last_updated) from wwv_flow_files where UPDATED_BY = p_user)
--         and id = (select max(id) from wwv_flow_files where updated_by = p_user);
         
    select blob_content into v_blob_data
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseSite;
         
    select filename into v_filename
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseSite;
         
END IF;
 
v_blob_len := dbms_lob.getlength(v_blob_data);
v_position := 1;

v_line_count := 0;

 
-- Read and convert binary to char
WHILE ( v_position <= v_blob_len ) LOOP
 v_raw_chunk := dbms_lob.substr(v_blob_data,c_chunk_len,v_position);
 v_char :=  chr(hex_to_decimal(rawtohex(v_raw_chunk)));
 v_line := v_line || v_char;
 v_position := v_position + c_chunk_len;

-- When a whole line is retrieved
 IF v_char = CHR(10) THEN

-- added for ignoring first row   
   v_line_count := v_line_count + 1;

-- Convert comma to : to use wwv_flow_utilities
   v_line := REPLACE (v_line, ',', ':');

-- Convert each column separated by : into array of data
   v_data_array := wwv_flow_utilities.string_to_table (v_line);



  IF v_line_count > 1 THEN

-- Insert data into target table
   
   EXECUTE IMMEDIATE 'insert into upload_sites (site_id, site_fid, site_year, fieldoffice_id, 
        field_office, project_id, project, 
        segment_id, segment, season_id, season, bend, bendrn, bend_river_mile, comments,
        last_updated, upload_session_id,uploaded_by,upload_filename)
        
         values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19)'
         USING 
         to_number(v_data_array(1)),
         v_data_array(2),
         to_number(v_data_array(3)),
         v_data_array(4),
         v_data_array(5),         
         to_number(v_data_array(6)),
         v_data_array(7),
         to_number(v_data_array(8)),         
         v_data_array(9),
         v_data_array(10), 
         v_data_array(11),
         to_number(v_data_array(12)),
         v_data_array(13),
         to_number(v_data_array(14)),
         v_data_array(15),
         sysdate,
         l_session,
         p_user,
         v_filename;

  END IF; -- v_line_count > 1
 
-- Clear out
   v_line := NULL;
   
 END IF;
  
 END LOOP;

  v_line_count := 0;

  -- Loop through upload_site records for this session
  FOR Y in (
    SELECT site_year||fieldoffice_id||project_id||segment_id||season_id||bend||bendrn as upload_fields
    FROM upload_sites
    WHERE uploaded_by = p_user
     and upload_session_id = l_session
     and site_fid IS NOT NULL)
  
  LOOP
  
    SELECT count(*)
        INTO l_siteMatch
    FROM ds_sites
    WHERE (year||fieldoffice||project_id||segment_id||season||bend||bendrn) = y.upload_fields;
    
    l_siteMatch_tot := l_siteMatch_tot + l_siteMatch; -- total sites that match 
  
  END LOOP;
  
  -- If there are uploaded sites that already exist in DS_SITES then don't upload any records to DS_SITES
  IF l_siteMatch_tot = 0 THEN
  
      -- insert into ds_sites - only those with site_FID
      INSERT into ds_sites (site_fid, year, fieldoffice, project_id,
            segment_id, season, bend, bendrn,
            last_updated, uploaded_by,upload_filename,upload_session_id)
        
        SELECT site_fid, site_year, fieldoffice_id, project_id,
            segment_id, season_id, bend, bendrn,
            last_updated, uploaded_by,upload_filename, upload_session_id
        FROM UPLOAD_SITES
        WHERE uploaded_by = p_user
         and upload_session_id = l_session
         and site_fid IS NOT NULL;
         
      COMMIT;
      
  ELSE
   
    -- number of uploaded sites that already exist in DS_SITES
    p_siteMatch := l_siteMatch_tot;
    
  END IF;
  
  -- count how many new sites were added to DS_SITES
  select count(*)
   into p_cnt
  from ds_sites
  where upload_session_id = l_session;
  
  

 
-- Delete file from wwv_flows
-- DELETE FROM wwv_flow_files
-- WHERE id = (select max(id) from wwv_flow_files where updated_by = p_user);
 
 DELETE from APEX_APPLICATION_FILES 
 where name = p_fileBrowseSite;
 
 
 
 l_session := 0;
 l_siteMatch := 0;
 l_siteMatch_tot := 0;
 
 
END uploadSiteDatasheet;


PROCEDURE uploadSiteDatasheetCheck (
                                p_user IN upload_sites.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_siteMatch OUT number,
                                p_fileBrowseSite IN varchar2 default null,
                                p_siteSessionID OUT number
                            ) IS

v_blob_data       BLOB;
v_blob_len        NUMBER;
v_position        NUMBER;
v_raw_chunk       RAW(10000);
v_char      CHAR(1);
c_chunk_len   number       := 1;
v_line        VARCHAR2 (32767)        := NULL;
v_data_array      wwv_flow_global.vc_arr2;

v_filename varchar2(200);

-- added for ignoring first row
v_line_count number;

l_session number default 0;
l_fish_max_session number default 0;
l_mr_max_session number default 0;
l_supp_max_session number default 0;

l_siteMatch number default 0;
l_siteMatch_tot number default 0;


BEGIN

-- set session number from sequence
l_session := site_load_session_seq.nextval;

p_siteSessionID := l_session;

p_siteMatch := l_siteMatch_tot;

-- Read data from wwv_flow_files
IF p_user is not null THEN
        
    select blob_content into v_blob_data
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseSite;
         
    select filename into v_filename
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseSite;
         
END IF;
 
v_blob_len := dbms_lob.getlength(v_blob_data);
v_position := 1;

v_line_count := 0;

 
-- Read and convert binary to char
WHILE ( v_position <= v_blob_len ) LOOP
 v_raw_chunk := dbms_lob.substr(v_blob_data,c_chunk_len,v_position);
 v_char :=  chr(hex_to_decimal(rawtohex(v_raw_chunk)));
 v_line := v_line || v_char;
 v_position := v_position + c_chunk_len;

-- When a whole line is retrieved
 IF v_char = CHR(10) THEN

-- added for ignoring first row   
   v_line_count := v_line_count + 1;

-- Convert comma to : to use wwv_flow_utilities
   v_line := REPLACE (v_line, ',', ':');

-- Convert each column separated by : into array of data
   v_data_array := wwv_flow_utilities.string_to_table (v_line);



  IF v_line_count > 1 THEN

-- Insert data into target table
   
   EXECUTE IMMEDIATE 'insert into upload_sites (site_id, site_fid, site_year, fieldoffice_id, 
        field_office, project_id, project, 
        segment_id, segment, season_id, season, bend, bendrn, bend_river_mile, comments,
        last_updated, upload_session_id,uploaded_by,upload_filename)
        
         values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19)'
         USING 
         to_number(v_data_array(1)),
         v_data_array(2),
         to_number(v_data_array(3)),
         v_data_array(4),
         v_data_array(5),         
         to_number(v_data_array(6)),
         v_data_array(7),
         to_number(v_data_array(8)),         
         v_data_array(9),
         v_data_array(10), 
         v_data_array(11),
         to_number(v_data_array(12)),
         v_data_array(13),
         to_number(v_data_array(14)),
         v_data_array(15),
         sysdate,
         l_session,
         p_user,
         v_filename;

  END IF; -- v_line_count > 1
 
-- Clear out
   v_line := NULL;
   
 END IF;
  
 END LOOP;

  v_line_count := 0;

  -- Loop through upload_site records for this session THAT ARE NEW SITES (site_fid is not null)
  FOR Y in (
    SELECT site_year||fieldoffice_id||project_id||segment_id||season_id||bend||bendrn as upload_fields
    FROM upload_sites
    WHERE uploaded_by = p_user
     and upload_session_id = l_session
     and site_fid IS NOT NULL)
  
  LOOP
  
    SELECT count(*)
        INTO l_siteMatch
    FROM ds_sites
    WHERE (year||fieldoffice||project_id||segment_id||season||bend||bendrn) = y.upload_fields;
    
    l_siteMatch_tot := l_siteMatch_tot + l_siteMatch; -- total sites that match 
  
  END LOOP;
  
  -- If there are NEW uploaded sites that already exist in DS_SITES then don't upload any records to DS_SITES
  IF l_siteMatch_tot = 0 THEN
  
      -- insert into ds_sites - only those with site_FID
      INSERT into ds_sites_check (site_fid, year, fieldoffice, project_id,
            segment_id, season, bend, bendrn,
            last_updated, uploaded_by,upload_filename,upload_session_id)
        
        SELECT site_fid, site_year, fieldoffice_id, project_id,
            segment_id, season_id, bend, bendrn,
            last_updated, uploaded_by,upload_filename, upload_session_id
        FROM UPLOAD_SITES
        WHERE uploaded_by = p_user
         and upload_session_id = l_session
         and site_fid IS NOT NULL;
         
      COMMIT;
      
  ELSE
   
    -- number of uploaded sites that already exist in DS_SITES
    p_siteMatch := l_siteMatch_tot;
    
  END IF;
  
  -- count how many new sites were added to DS_SITES
  -- Not Needed now - 5/13/15 - JRF
--  select count(*)
--   into p_cnt
--  from ds_sites_check
--  where upload_session_id = l_session;
  
  

 
-- Delete file from wwv_flows
-- DELETE FROM wwv_flow_files
-- WHERE id = (select max(id) from wwv_flow_files where updated_by = p_user);
 
 DELETE from APEX_APPLICATION_FILES 
 where name = p_fileBrowseSite; 
 
 l_session := 0;
 l_siteMatch := 0;
 l_siteMatch_tot := 0;
 
 
END uploadSiteDatasheetCheck;


PROCEDURE uploadSiteDatasheetCheck2020 (
                                p_user IN upload_sites.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_siteMatch OUT number,
                                p_fileBrowseSite IN varchar2 default null,
                                p_siteSessionID OUT number
                            ) IS

v_blob_data       BLOB;
v_blob_len        NUMBER;
v_position        NUMBER;
v_raw_chunk       RAW(10000);
v_char      CHAR(1);
c_chunk_len   number       := 1;
v_line        VARCHAR2 (32767)        := NULL;
v_data_array      wwv_flow_global.vc_arr2;

v_filename varchar2(200);

-- added for ignoring first row
v_line_count number;

l_session number default 0;
l_fish_max_session number default 0;
l_mr_max_session number default 0;
l_supp_max_session number default 0;

l_siteMatch number default 0;
l_siteMatch_tot number default 0;


BEGIN

-- set session number from sequence
l_session := site_load_session_seq.nextval;

p_siteSessionID := l_session;

p_siteMatch := l_siteMatch_tot;

-- Read data from wwv_flow_files
IF p_user is not null THEN
        
    select blob_content into v_blob_data
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseSite;
         
    select filename into v_filename
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseSite;
         
END IF;
 
v_blob_len := dbms_lob.getlength(v_blob_data);
v_position := 1;

v_line_count := 0;

 
-- Read and convert binary to char
WHILE ( v_position <= v_blob_len ) LOOP
 v_raw_chunk := dbms_lob.substr(v_blob_data,c_chunk_len,v_position);
 v_char :=  chr(hex_to_decimal(rawtohex(v_raw_chunk)));
 v_line := v_line || v_char;
 v_position := v_position + c_chunk_len;

-- When a whole line is retrieved
 IF v_char = CHR(10) THEN

-- added for ignoring first row   
   v_line_count := v_line_count + 1;

-- Convert comma to : to use wwv_flow_utilities
   v_line := REPLACE (v_line, ',', ':');

-- Convert each column separated by : into array of data
   v_data_array := wwv_flow_utilities.string_to_table (v_line);



  IF v_line_count > 1 THEN

-- Insert data into target table
--This code is a challenge to understand what is going on. No idea how the data is supposed to flow.
   --2020 Corals latest email  has
  -- A= site_id || b = site_fid || c = SITE_YEAR || d = FIELDOFFICE_ID || e = FIELD_OFFICE ||
  -- F = PROJECT_ID || G= 	PROJECT|| H=SEGMENT_ID	 || I = SEGMENT || J=	SEASON_ID	|| K = SEASON	||
  -- L= BEND	|| M = BENDRN	 || N= BEND_RIVER_MILE	||O =COMMENTS


   
   EXECUTE IMMEDIATE 'insert into upload_sites (site_id, site_fid, site_year, fieldoffice_id, 
        field_office, project_id, project, 
        segment_id, segment, season_id, season, bend, bendrn, bend_river_mile, comments,
        last_updated, upload_session_id,uploaded_by,upload_filename)
        
         values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19)'
         USING 
         to_number(v_data_array(1)),
         v_data_array(2),
         to_number(v_data_array(3)),
         v_data_array(4),
         v_data_array(5),         
         to_number(v_data_array(6)),
         v_data_array(7),
         to_number(v_data_array(8)),         
         v_data_array(9),
         v_data_array(10), 
         v_data_array(11),
         to_number(v_data_array(12)),
         v_data_array(13),
         to_number(v_data_array(14)),
         v_data_array(15),
         sysdate,
         l_session,
         p_user,
         v_filename;

  END IF; -- v_line_count > 1
 
-- Clear out
   v_line := NULL;
   
 END IF;
  
 END LOOP;

  v_line_count := 0;

  -- Loop through upload_site records for this session THAT ARE NEW SITES (site_fid is not null)
  FOR Y in (
    SELECT site_year||fieldoffice_id||project_id||segment_id||season_id||bend||bendrn as upload_fields
    FROM upload_sites
    WHERE uploaded_by = p_user
     and upload_session_id = l_session
     and site_fid IS NOT NULL)
  
  LOOP
  
    SELECT count(*)
       INTO l_siteMatch
    FROM ds_sites
    WHERE (year||fieldoffice||project_id||segment_id||season||bend||bendrn) = y.upload_fields;
    
    l_siteMatch_tot := l_siteMatch_tot + l_siteMatch; -- total sites that match 
  
  END LOOP;
  
  -- If there are NEW uploaded sites that already exist in DS_SITES then don't upload any records to DS_SITES
  IF l_siteMatch_tot = 0 THEN
  
      -- insert into ds_sites - only those with site_FID
      INSERT into ds_sites_check (site_fid, year, fieldoffice, project_id,
            segment_id, season, bend, bendrn,
            last_updated, uploaded_by,upload_filename,upload_session_id)
        
        SELECT site_fid, site_year, fieldoffice_id, project_id,
            segment_id, season_id, bend, bendrn,
            last_updated, uploaded_by,upload_filename, upload_session_id
        FROM UPLOAD_SITES
        WHERE uploaded_by = p_user
         and upload_session_id = l_session
         and site_fid IS NOT NULL;
         
      COMMIT;
      
  ELSE
   
    -- number of uploaded sites that already exist in DS_SITES
    p_siteMatch := l_siteMatch_tot;
    
  END IF;
  
  -- count how many new sites were added to DS_SITES
  -- Not Needed now - 5/13/15 - JRF
--  select count(*)
--   into p_cnt
--  from ds_sites_check
--  where upload_session_id = l_session;
  
  

 
-- Delete file from wwv_flows
-- DELETE FROM wwv_flow_files
-- WHERE id = (select max(id) from wwv_flow_files where updated_by = p_user);
 
 DELETE from APEX_APPLICATION_FILES 
 where name = p_fileBrowseSite; 
 
 l_session := 0;
 l_siteMatch := 0;
 l_siteMatch_tot := 0;
 
 
END uploadSiteDatasheetCheck2020;

------------

PROCEDURE uploadMRdatasheetCheck (
                                p_user IN upload_mr.uploaded_by%TYPE,
                                -- p_complete IN upload_mr.complete%TYPE,
                                p_checkby IN upload_mr.checkby%TYPE,
                                p_cnt OUT number,
                                p_mrfidMatch OUT number,
                                p_fileBrowseMR IN varchar2 default null,                                
                                p_mrSessionID OUT number
                            ) IS

v_blob_data       BLOB;
v_blob_len        NUMBER;
v_position        NUMBER;
v_raw_chunk       RAW(10000);
v_char      CHAR(1);
c_chunk_len   number := 1;
v_line        VARCHAR2 (32767) := NULL;
v_data_array      wwv_flow_global.vc_arr2;

v_filename varchar2(200);

-- added for ignoring first row
v_line_count number;

l_session number default 0;
l_max_mr_id number default 0;
-- l_site_id number default 0;
l_mr_max_session number default 0;
l_site_max_session number default 0;

l_mrDuplicate number default 0;
l_mrDuplicateTot number default 0;


BEGIN

-- set session number from sequence
l_session := mr_load_session_seq.nextval;

p_mrSessionID := l_session;

p_cnt := 0;
p_mrfidMatch := 0;

-- test

-- Read data from wwv_flow_files
 IF p_user is not null THEN
         
    select blob_content into v_blob_data
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseMR;
         
    select filename into v_filename
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseMR;

 END IF;

v_blob_len := dbms_lob.getlength(v_blob_data);
v_position := 1;

v_line_count := 0;

 
-- Read and convert binary to char
WHILE ( v_position <= v_blob_len ) LOOP
 v_raw_chunk := dbms_lob.substr(v_blob_data,c_chunk_len,v_position);
 v_char :=  chr(hex_to_decimal(rawtohex(v_raw_chunk)));
 v_line := v_line || v_char;
 v_position := v_position + c_chunk_len;

-- When a whole line is retrieved
     IF v_char = CHR(10) THEN

    -- added for ignoring first row   
       v_line_count := v_line_count + 1;

    -- Convert comma to : to use wwv_flow_utilities
       v_line := REPLACE (v_line, ':', '_');
       v_line := REPLACE (v_line, ',', ':');

    -- Convert each column separated by : into array of data
       v_data_array := wwv_flow_utilities.string_to_table (v_line);



       IF v_line_count > 1 THEN

    -- Insert data into target table
       EXECUTE IMMEDIATE 'insert into upload_mr (site_id, site_fid, mr_fid, season, setdate, subsample, subsamplepass, 
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
        uploaded_by, upload_filename, checkby,
        no_turbidity, no_velocity)

         values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23,:24,
            :25,:26,:27,:28,:29,:30,:31,:32,:33,:34,:35,:36,:37,:38,:39,:40,:41,:42,:43,:44,:45,:46,:47,
            :48,:49,:50,:51,:52,:53,:54,:55,:56,:57,:58,:59,:60,:61,:62,:63,:64,:65,:66,:67,:68,:69,:70,
            :71)'
         USING 
         to_number(v_data_array(1)),
         v_data_array(2),
         v_data_array(3),
         v_data_array(4),
         to_date(v_data_array(5),'mm/dd/yyyy'),
         to_number(v_data_array(6)),
         to_number(v_data_array(7)),
         v_data_array(8),
         v_data_array(9),
         
         v_data_array(10),
         v_data_array(11),
         to_number(v_data_array(12)),
         to_number(v_data_array(13)),
         to_number(v_data_array(14)),
         to_number(v_data_array(15)),
         
         to_number(v_data_array(16)),
         to_number(v_data_array(17)),
         
         to_number(v_data_array(18)),
         v_data_array(19),
         
         v_data_array(20),
         to_number(v_data_array(21)),
         to_number(v_data_array(22)),
         
         v_data_array(23),
         v_data_array(24),
         v_data_array(25),
         v_data_array(26),
         v_data_array(27),
         v_data_array(28),
         v_data_array(29),
         
         v_data_array(30),
         v_data_array(31),
         v_data_array(32),
         v_data_array(33),
         
         v_data_array(34),
         v_data_array(35),
         v_data_array(36),
         v_data_array(37),
         v_data_array(38),
         v_data_array(39),
         
         v_data_array(40),
         to_number(v_data_array(41)),
         to_number(v_data_array(42)),
         v_data_array(43),
         to_number(v_data_array(44)),
         to_number(v_data_array(45)),
         
         to_number(v_data_array(46)),
         to_number(v_data_array(47)),
         to_number(v_data_array(48)),
         to_number(v_data_array(49)),
         
         to_number(v_data_array(50)),
         to_number(v_data_array(51)),
         to_number(v_data_array(52)),
         to_number(v_data_array(53)),
         
         to_number(v_data_array(54)),
         to_number(v_data_array(55)),
         to_number(v_data_array(56)),
         to_number(v_data_array(57)),
         
         to_number(v_data_array(58)),
         to_number(v_data_array(59)),
         to_number(v_data_array(60)),
         to_number(v_data_array(61)),
         to_number(v_data_array(62)),
         to_number(v_data_array(63)),         
         
         v_data_array(66),  --  v_data_array(66),
         sysdate,
         l_session,
         p_user,
         v_filename,
         p_checkby,
         v_data_array(64),
         v_data_array(65);
         
         commit;

       END IF; -- v_line_count > 1
         
        -- Clear out
           v_line := NULL;
           
     END IF;
  
  END LOOP;

  v_line_count := 0; 
  
  --  get latest mr_id
  SELECT distinct max(mr_id)
    into l_max_mr_id
  FROM upload_mr;  

     
  -- Loop through upload_mr records for this session
  FOR Y in (
    SELECT mr_fid
    FROM upload_mr
    WHERE uploaded_by = p_user
     and upload_session_id = l_session)
  
  LOOP
  
    SELECT count(*)
        INTO l_mrDuplicate
    FROM ds_moriver
    WHERE mr_fid = y.mr_fid;
    
    l_mrDuplicateTot := l_mrDuplicateTot +l_mrDuplicate; -- total mr_fid that match 
  
  END LOOP;
  
  -- number of mr_fid that already exist in DS_MORIVER
    p_mrfidMatch := l_mrDuplicateTot;
  
  -- If there are mr_fid that already exist in DS_MORIVER then don't upload any records to DS_MORIVER  
  
      -- insert into ds_moriver
      INSERT into ds_moriver_check (site_id, mr_fid, season, setdate, subsample, subsamplepass, comments, last_updated,
            uploaded_by, subsamplen, recorder, gear, gear_type, temp, turbidity, conductivity, do,
            distance, width, netrivermile, structurenumber, usgs, riverstage, discharge,
            u1, u2, u3, u4, u5, u6, u7, macro, meso, habitatrn, qc,
            micro_structure, structure_flow, structure_mod, set_site_1, set_site_2, set_site_3,
            micro,
            starttime, startlatitude, startlongitude, stoptime, stoplatitude, stoplongitude,
            depth1, velocitybot1, velocity08_1, velocity02or06_1,
            depth2, velocitybot2, velocity08_2, velocity02or06_2,
            depth3, velocitybot3, velocity08_3, velocity02or06_3,
            watervel, cobble, organic, silt, sand, gravel, upload_filename,
            complete, checkby, 
            no_turbidity, no_velocity, 
            upload_session_id, site_fid
            )
        SELECT site_id, mr_fid, season, setdate, subsample, subsamplepass, comments, last_updated, uploaded_by,
            subsamplen, recorder, gear, gear_type, temp, turbidity, conductivity, do,
            distance, width, netrivermile, structurenumber, usgs, riverstage, discharge,
            u1, u2, u3, u4, u5, u6, u7, macro, meso, habitatrn, qc,
            micro_structure, structure_flow, structure_mod, set_site_1, set_site_2, set_site_3,
            micro_structure || structure_flow || structure_mod || set_site_1 || set_site_2 || set_site_3, 
            replace(starttime,'_',':') as start_time, startlatitude, startlongitude, 
            replace(stoptime,'_',':') as stop_time, stoplatitude, stoplongitude,
            depth1, velocitybot1, velocity08_1, velocity02or06_1,
            depth2, velocitybot2, velocity08_2, velocity02or06_2,
            depth3, velocitybot3, velocity08_3, velocity02or06_3,
            watervel, cobble, organic, silt, sand, gravel, upload_filename,
            complete, checkby, 
            no_turbidity, no_velocity, 
            upload_session_id, site_fid
        FROM UPLOAD_MR
        WHERE uploaded_by = p_user
         and upload_session_id = l_session
         and NOT EXISTS (select mr_fid from ds_moriver
                         where ds_moriver.mr_fid = upload_mr.mr_fid);
         
      COMMIT; 

  
  -- count how many new mr records were added to DS_MORIVER
  -- Not Needed now - 5/13/15 - JRF
--  select count(*)
--   into p_cnt
--  from ds_moriver_check
--  where upload_session_id = l_session;
  
  
   -- clear out variables
    l_site_max_session := 0;
    l_mr_max_session := 0;
    
    l_mrDuplicate := 0;
    l_mrDuplicateTot := 0;
 
 
-- Delete file from wwv_flows
-- DELETE FROM wwv_flow_files
-- WHERE id = (select max(id) from wwv_flow_files where updated_by = p_user);
 
 DELETE FROM APEX_APPLICATION_FILES 
 where name = p_fileBrowseMR;


 l_session := 0;
 
END uploadMRdatasheetCheck;

-----

PROCEDURE uploadFishDatasheetCheck (
                                p_user IN upload_fish.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_ffidMatch OUT number,
                                p_fileBrowseF IN varchar2 default null,
                                p_fishSessionID OUT number,
                                p_mrSessionID in number
                            ) IS

v_blob_data       BLOB;
v_blob_len        NUMBER;
v_position        NUMBER;
v_raw_chunk       RAW(10000);
v_char      CHAR(1);
c_chunk_len   number := 1;
v_line        VARCHAR2 (32767) := NULL;
v_data_array      wwv_flow_global.vc_arr2;

v_filename varchar2(200);

-- added for ignoring first row
v_line_count number;

l_session number default 0;
 l_fish_max_session number default 0;
 l_mr_max_session number default 0;
 
l_fDuplicate number default 0;
l_fDuplicateTot number default 0;


BEGIN

-- set session number from sequence
-- l_session := mr_load_session_seq.nextval;

-- l_session := fish_load_session_seq.nextval;  -- Changed to use MR Session for all so MR, Fish and Supp have same Session ID - 11/21/13

l_session := p_mrSessionID;  -- New 11/21/13

p_fishSessionID := l_session;

p_cnt := 0;
p_ffidMatch := 0;

-- Read data from wwv_flow_files
 IF p_user is not null THEN
 
        
    select blob_content into v_blob_data
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseF;
         
    select filename into v_filename
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseF;       
    
     
 END IF;  
     
     
 
v_blob_len := dbms_lob.getlength(v_blob_data);
v_position := 1;

v_line_count := 0;

 
-- Read and convert binary to char
WHILE ( v_position <= v_blob_len ) LOOP
 v_raw_chunk := dbms_lob.substr(v_blob_data,c_chunk_len,v_position);
 v_char :=  chr(hex_to_decimal(rawtohex(v_raw_chunk)));
 v_line := v_line || v_char;
 v_position := v_position + c_chunk_len;

-- When a whole line is retrieved
     IF v_char = CHR(10) THEN

    -- added for ignoring first row   
       v_line_count := v_line_count + 1;

    -- Convert comma to : to use wwv_flow_utilities
       v_line := REPLACE (v_line, ',', ':');

    -- Convert each column separated by : into array of data
       v_data_array := wwv_flow_utilities.string_to_table (v_line);



       IF v_line_count > 1 THEN

    -- Insert data into target table
       EXECUTE IMMEDIATE 'insert into upload_fish (site_id, f_fid, mr_fid, panelhook, bait, species, length, weight,
        fishcount, fin_curl, otolith, rayspine, scale, ftprefix, ftnum, ftmr,
        comments, last_updated, upload_session_id,uploaded_by, upload_filename)
         values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21)'
         USING 
         to_number(v_data_array(1)),
         v_data_array(2),
         v_data_array(3),
         v_data_array(4),
         v_data_array(5),
         v_data_array(6),
         to_number(v_data_array(7)),
         to_number(v_data_array(8)),
         
         to_number(v_data_array(9)),
         v_data_array(10),
         v_data_array(11),
         v_data_array(12),
         v_data_array(13),
         v_data_array(14),
         v_data_array(15),
         v_data_array(16),         
         v_data_array(17),
         sysdate,
         l_session,
         p_user,
         v_filename;
         
         commit;

       END IF; -- v_line_count > 1
         
        -- Clear out
           v_line := NULL;
           
     END IF;
  
  END LOOP;

  v_line_count := 0; 
  
 -- Loop through upload_fish records for this session
  FOR Y in (
    SELECT f_fid
    FROM upload_fish
    WHERE uploaded_by = p_user
     and upload_session_id = l_session)
  
  LOOP
  
    SELECT count(*)
        INTO l_fDuplicate
    FROM ds_fish
    WHERE f_fid = y.f_fid;
    
    l_fDuplicateTot := l_fDuplicateTot +l_fDuplicate; -- total f_fid that match 
  
  END LOOP;
  
  -- number of mr_fid that already exist in DS_MORIVER
    p_ffidMatch := l_fDuplicateTot;
  
  -- If there are f_fid that already exist in DS_FISH then don't upload any records to DS_FISH  
   
      -- insert into ds_fish
      INSERT into ds_fish_check (mr_id, f_fid, panelhook, bait, species, length, weight,
        last_updated, uploaded_by,fishcount, fin_curl, otolith, rayspine, scale, ftprefix, ftnum, ftmr,
        upload_filename, complete, checkby, upload_session_id, mr_fid)
        SELECT mr_id, f_fid, panelhook, bait, species, length, weight, last_updated, uploaded_by,
            fishcount, fin_curl, otolith, rayspine, scale, ftprefix, ftnum, ftmr, upload_filename,
            complete, checkby, upload_session_id, mr_fid
        FROM UPLOAD_FISH
        WHERE uploaded_by = p_user
         and upload_session_id = l_session
         and species IS NOT NULL  -- added to remove entries where no species was entered 
         and NOT EXISTS (select f_fid from ds_fish
                         where ds_fish.f_fid = upload_fish.f_fid); -- added to keep duplicates from being entered         
         
      COMMIT;
  
  
  
  -- count how many new mr records were added to DS_MORIVER
  -- Not Needed now - 5/13/15 - JRF
--  select count(*)
--   into p_cnt
--  from ds_fish_check
--  where upload_session_id = l_session; 
  
 
 
-- clear out variables
l_fish_max_session := 0;
l_mr_max_session := 0;

l_fDuplicate := 0;
l_fDuplicateTot := 0;
 
-- Delete file from wwv_flows
-- DELETE FROM wwv_flow_files
-- WHERE id = (select max(id) from wwv_flow_files where updated_by = p_user);
 
 DELETE from APEX_APPLICATION_FILES 
 where name = p_fileBrowseF;
 
 


 l_session := 0;
 
END uploadFishDatasheetCheck;

-----------
PROCEDURE uploadSuppDatasheetCheck (
                                p_user IN upload_supplemental.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_sfidMatch OUT number,
                                p_fileBrowseSup IN varchar2 default null,
                                p_suppSessionID OUT number,
                                p_mrSessionID in number,
                                p_plus_cnt OUT number,
                                p_plus_FID OUT varchar2
                            ) IS

v_blob_data       BLOB;
v_blob_len        NUMBER;
v_position        NUMBER;
v_raw_chunk       RAW(10000);
v_char      CHAR(1);
c_chunk_len   number       := 1;
v_line        VARCHAR2 (32767)        := NULL;
v_data_array      wwv_flow_global.vc_arr2;

v_filename varchar2(200);

-- added for ignoring first row
v_line_count number;

l_session number default 0;
l_fish_max_session number default 0;
l_mr_max_session number default 0;
l_supp_max_session number default 0;

l_sDuplicate number default 0;
l_sDuplicateTot number default 0;

l_plusCount number default 0;
l_plusCountTot number default 0;
l_f_fid varchar(100);
l_f_fidTot varchar(1000);

BEGIN

-- set session number from sequence
-- l_session := supp_load_session_seq.nextval;  -- Changed to use MR Session for all so MR, Fish and Supp have same Session ID - 11/21/13

l_session := p_mrSessionID;  -- New 11/21/13

p_suppSessionID := l_session;

p_cnt := 0;
p_sfidMatch := 0;

-- Read data from wwv_flow_files
 IF p_user is not null THEN
             
    select blob_content into v_blob_data
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseSup;
         
    select filename into v_filename
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseSup;
             
 END IF;  

 
v_blob_len := dbms_lob.getlength(v_blob_data);
v_position := 1;

v_line_count := 0;

 
-- Read and convert binary to char
WHILE ( v_position <= v_blob_len ) LOOP
 v_raw_chunk := dbms_lob.substr(v_blob_data,c_chunk_len,v_position);
 v_char :=  chr(hex_to_decimal(rawtohex(v_raw_chunk)));
 v_line := v_line || v_char;
 v_position := v_position + c_chunk_len;

-- When a whole line is retrieved
 IF v_char = CHR(10) THEN

-- added for ignoring first row   
   v_line_count := v_line_count + 1;

-- Convert comma to : to use wwv_flow_utilities
   v_line := REPLACE (v_line, ',', ':');

-- Convert each column separated by : into array of data
   v_data_array := wwv_flow_utilities.string_to_table (v_line);



  IF v_line_count > 1 THEN

-- Insert data into target table
   EXECUTE IMMEDIATE 'insert into upload_supplemental (site_id, f_fid, mr_fid, 
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
            :21,:22,:23,:24,:25,:26,:27,:28,:29,:30,:31,:32,:33,:34,:35,:36,:37,:38,:39,:40,:41,:42,:43,:44,:45)'
            
         USING 
         to_number(v_data_array(1)),
         -- v_data_array(2),
         v_data_array(2),
         v_data_array(3),
         v_data_array(4),         
         v_data_array(5),
         v_data_array(6),
         to_number(v_data_array(7)),         
         v_data_array(8),
         to_number(v_data_array(9)),         
         v_data_array(10),
         v_data_array(11),
         v_data_array(12),
         v_data_array(13),
         v_data_array(14),
         v_data_array(15),         
         v_data_array(16),
         v_data_array(17),
         to_number(v_data_array(18)),
         to_number(v_data_array(19)),
         to_number(v_data_array(20)),         
         to_number(v_data_array(21)),
         to_number(v_data_array(22)),
         to_number(v_data_array(23)),
         to_number(v_data_array(24)),
         to_number(v_data_array(25)),
         to_number(v_data_array(26)),         
         to_number(v_data_array(27)),
         to_number(v_data_array(28)),
         to_number(v_data_array(29)),
         to_number(v_data_array(30)),
         to_number(v_data_array(31)),
         to_number(v_data_array(32)), 
         v_data_array(33), -- status
         v_data_array(34), -- hatchery origin
         v_data_array(35), 
         v_data_array(36), 
         v_data_array(37), 
         v_data_array(38),          
         v_data_array(39),
         v_data_array(40),
         v_data_array(41),    
         sysdate,
         l_session,
         p_user,
         v_filename;

  END IF; -- v_line_count > 1
 
-- Clear out
   v_line := NULL;
   
 END IF;
  
 END LOOP;

  v_line_count := 0;
 
 
 -- Loop through upload_supplemental records for this session
  FOR Y in (
    SELECT f_fid
    FROM upload_supplemental
    WHERE uploaded_by = p_user
     and upload_session_id = l_session)
  
  LOOP
  
    SELECT count(*)
        INTO l_sDuplicate
    FROM ds_supplemental
    WHERE f_fid = y.f_fid;
    
    l_sDuplicateTot := l_sDuplicateTot + l_sDuplicate; -- total f_fid that match 
  
  END LOOP;
  
  -- number of mr_fid that already exist in DS_MORIVER
    p_sfidMatch := l_sDuplicateTot;  
 
  -- If there are f_fid that already exist in DS_FISH then don't upload any records to DS_FISH 

      -- insert into ds_supplemental
      INSERT into ds_supplemental_check (mr_id, f_id, f_fid, tagnumber, pitrn, cwtyn, dangler,
        last_updated, uploaded_by, scuteloc, scutenum,
        elhv, elcolor, erhv, ercolor, genetic, genetics_vial_number,
        head, snouttomouth, inter, mouthwidth, m_ib, l_ob, l_ib, r_ib, r_ob, anal, dorsal,
        sex, stage, status, hatchery_origin, recapture, photo,
        genetic_needs, broodstock, hatch_wild, species_id, archive,
        other_tag_info,comments, upload_filename, complete, checkby, upload_session_id,
        mr_fid, scuteloc2, scutenum2)
        SELECT mr_id, f_id, f_fid, tagnumber, pitrn, cwtyn, dangler, last_updated, uploaded_by,
            scuteloc, scutenum, elhv, elcolor, erhv, ercolor, genetic, genetics_vial_number,
            head, snouttomouth, inter, mouthwidth, m_ib, l_ob, l_ib, r_ib, r_ob, anal, dorsal,
            sex, stage, status, hatchery_origin, recapture, photo,
            genetic_needs, broodstock, hatch_wild, species_id, archive,
            other_tag_info, comments, upload_filename, complete, checkby,upload_session_id,
            mr_fid, scuteloc2, scutenum2
        FROM UPLOAD_SUPPLEMENTAL
        WHERE uploaded_by = p_user
         and upload_session_id = l_session
         and NOT EXISTS (select f_fid from ds_supplemental
                         where ds_supplemental.f_fid = upload_supplemental.f_fid); -- added to keep duplicates from being entered
         
      COMMIT;
      
      
  -- NEW SECTION -- To catch '+' in Tag number for Supplemental - 1/6/14
  -- Loop through upload_supplemental records for this session
  FOR Z in (
    SELECT f_fid, tagnumber
    FROM ds_supplemental_check
    WHERE uploaded_by = p_user
     and upload_session_id = l_session)
  
  LOOP
  
    IF LENGTH(z.tagnumber) - LENGTH(replace(z.tagnumber,'+','')) != 0 THEN -- check for specific character - '+'
    
        l_plusCount := 1;
        
        l_f_fid := z.f_fid;
    
    
    END IF;  
       
    l_plusCountTot := l_plusCountTot + l_plusCount; -- total tagnumber with + in string     
    
    
    IF l_f_fid IS NOT NULL THEN -- it got value in IF statement above
        l_f_fidTot := l_f_fidTot || ', ' || l_f_fid;
    END IF;
    
    -- Clear values
    l_plusCount := 0;
    l_f_fid := '';
  
  END LOOP;
  
  -- pass number of records with + in tagnumber back to APEX procedure
  p_plus_cnt := l_plusCountTot;
  
  
  p_plus_FID := ltrim(l_f_fidTot, ', ');   
  
  -----------------
      
  
  
  -- count how many new mr records were added to DS_MORIVER
  -- Not Needed now - 5/13/15 - JRF
--  select count(*)
--   into p_cnt
--  from ds_supplemental_check
--  where upload_session_id = l_session; 

-- clear out variables
l_fish_max_session := 0;
l_mr_max_session := 0;
l_supp_max_session := 0;

l_sDuplicate := 0;
l_sDuplicateTot := 0;

l_plusCountTot := 0;
l_f_fidTot := '';

 
-- Delete file from wwv_flows
-- DELETE FROM wwv_flow_files
-- WHERE id = (select max(id) from wwv_flow_files where updated_by = p_user);
 
 DELETE from APEX_APPLICATION_FILES 
 where name = p_fileBrowseSup; 
 
 l_session := 0;
 
 
END uploadSuppDatasheetCheck;

PROCEDURE uploadFinal (
                                p_user IN ds_sites_check.uploaded_by%TYPE,
                                p_site_cnt_final OUT number,
                                p_siteSessionID in number,
                                p_mrSessionID in number,
                                p_mr_cnt_final OUT number,
                                p_fishSessionID in number,
                                p_fishCntFinal OUT number,
                                p_suppSessionID in number,
                                p_suppCntFinal OUT number,
                                p_plus_cnt in number,
                                p_noSite_cnt OUT number,
                                p_noSiteID_msg OUT varchar2
                            ) IS
                            
-- Exceptions
 noSiteID exception;


  BEGIN 
    
    -- Set final counts to default of 0
    p_site_cnt_final := 0;
    p_mr_cnt_final := 0;
    p_fishCntFinal := 0;
    p_suppCntFinal := 0;   
  
   
  -- Sites
  
      -- insert into ds_sites - only those with site_FID
      INSERT into ds_sites (site_fid, year, fieldoffice, project_id,
            segment_id, season, bend, bendrn,
            last_updated, uploaded_by,upload_filename,upload_session_id, sample_unit_type)
        
        SELECT site_fid, year, fieldoffice, project_id,
            segment_id, season, bend, bendrn,
            last_updated, uploaded_by,upload_filename, upload_session_id, 'B'
        FROM ds_sites_check
        WHERE uploaded_by = p_user
         and upload_session_id = p_siteSessionID
         and site_fid IS NOT NULL;
         
      COMMIT;
      
  
  
  -- count how many new sites were added to DS_SITES
  select count(*)
   into p_site_cnt_final
  from ds_sites
  where upload_session_id = p_siteSessionID;
  
  -- Missouri River Data Sheets  
           
            
    -- if uploading a MR data sheet associated with new site
    FOR x IN
        (SELECT SITE_ID
         FROM DS_MORIVER_CHECK
         WHERE upload_session_id = p_mrSessionID)
     LOOP
        
        IF x.SITE_ID IS NULL OR x.SITE_ID = 0 THEN
        
             update DS_MORIVER_CHECK
             set site_id = 
                 (SELECT site_id
                  FROM DS_SITES
                  WHERE site_fid = ds_moriver_check.site_fid
                    and upload_session_id = p_siteSessionID
                  )
             WHERE upload_session_id = p_mrSessionID
                  AND site_id = x.site_id;
        END IF;
     
     END LOOP; 
     
     -- count how many records in DS_MORIVER_CHECK don't have SITE_ID after process above.
     -- If some records don't have SITE_ID then they tried to upload MR data with field-added site without uploading new sites
     SELECT count(*)
        INTO p_noSite_cnt     
     FROM  DS_MORIVER_CHECK
     WHERE SITE_ID IS NULL
     and upload_session_id = p_mrSessionID;
     
     -- Exception - if p_noSite_cnt is > 0 then do exception and jump to bottom
     IF NVL(p_noSite_cnt,0) > 0 THEN
        raise noSiteID;
     END IF;

     
     -- update DS_MORIVER_CHECK with COMPLETE field per MR_FID based on count of BAFI per MR_FID
     -- 11/12/14
     FOR X IN
        (SELECT MR_FID
         FROM DS_MORIVER_CHECK
         WHERE upload_session_id = p_mrSessionID)  
     LOOP
     
        FOR Y IN
            (SELECT count(species) as speciesCount
             FROM DS_FISH_CHECK
             WHERE upload_session_id = p_mrSessionID
             and species = 'BAFI'
             and mr_fid = x.mr_fid)        
        LOOP
        
            -- Commented out - 2/20/15 - not using COMPLETE any more.
            
--            IF Y.speciesCount is null or Y.speciesCount = 0 THEN
--            
--                UPDATE DS_MORIVER_CHECK
--                set COMPLETE = 1
--                WHERE upload_session_id = p_mrSessionID
--                 and mr_fid = x.mr_fid;
--             
--            END IF;  
            
            -- New section - 02/20//15 - If BAFI exists then blank out CHECKBY.  CHeckby already in DS_MORIVER_CHECK otherwise.
            IF Y.speciescount > 0 THEN
            
                UPDATE DS_MORIVER_CHECK
                set checkby = ''
                WHERE upload_session_id = p_mrSessionID
                 and mr_fid = x.mr_fid;            
            
            END IF;                  
        
        END LOOP;     
     
     END LOOP;     
     
     -- If there were any + symbols in tagnumber from Supplemental then don't insert data into ds_moriver
     IF p_plus_cnt < 1 or p_plus_cnt IS NULL THEN
     
     ------  insert records from temporary table to DS_MORIVER
     -- insert into ds_moriver
      INSERT into ds_moriver (site_id, mr_fid, season, setdate, subsample, subsamplepass, comments, last_updated,
            uploaded_by, subsamplen, recorder, gear, gear_type, temp, turbidity, conductivity, do,
            distance, width, netrivermile, structurenumber, usgs, riverstage, discharge,
            u1, u2, u3, u4, u5, u6, u7, macro, meso, habitatrn, qc,
            micro_structure, structure_flow, structure_mod, set_site_1, set_site_2, set_site_3,
            micro,
            starttime, startlatitude, startlongitude, stoptime, stoplatitude, stoplongitude,
            depth1, velocitybot1, velocity08_1, velocity02or06_1,
            depth2, velocitybot2, velocity08_2, velocity02or06_2,
            depth3, velocitybot3, velocity08_3, velocity02or06_3,
            watervel, cobble, organic, silt, sand, gravel, upload_filename,
            complete, checkby, 
            no_turbidity, no_velocity, 
            upload_session_id
            )
        SELECT site_id, mr_fid, season, setdate, subsample, subsamplepass, comments, last_updated, uploaded_by,
            subsamplen, recorder, gear, gear_type, temp, turbidity, conductivity, do,
            distance, width, netrivermile, structurenumber, usgs, riverstage, discharge,
            u1, u2, u3, u4, u5, u6, u7, macro, meso, habitatrn, qc,
            micro_structure, structure_flow, structure_mod, set_site_1, set_site_2, set_site_3,
            micro, 
            starttime, startlatitude, startlongitude, 
            stoptime, stoplatitude, stoplongitude,
            depth1, velocitybot1, velocity08_1, velocity02or06_1,
            depth2, velocitybot2, velocity08_2, velocity02or06_2,
            depth3, velocitybot3, velocity08_3, velocity02or06_3,
            watervel, cobble, organic, silt, sand, gravel, upload_filename,
            complete, checkby, 
            no_turbidity, no_velocity, 
            upload_session_id
        FROM DS_MORIVER_CHECK
        WHERE uploaded_by = p_user
         and upload_session_id = p_mrSessionID;
         
      COMMIT;     
     
      
     END IF;
      
  -- count how many new mr records were added to DS_MORIVER
  select count(*)
   into p_mr_cnt_final
  from ds_moriver
  where upload_session_id = p_mrSessionID;
  
  ----------------------- FISH  -------------------
  

 update ds_fish_check
 set (mr_id, complete, checkby) = 
     (SELECT mr_id, complete, checkby 
     FROM ds_moriver
     WHERE mr_fid = ds_fish_check.MR_FID
      and UPLOAD_SESSION_ID = p_mrSessionID )   
 WHERE upload_session_id = p_fishSessionID;
 
 -- If adding fish that weren't in original upload...
 FOR Y in (
    select mr_id, complete, checkby
    from ds_fish_check
    where upload_session_id = p_fishSessionID)
  
  LOOP
  
    IF y.mr_id IS NULL THEN
    
     update ds_fish_check
 			set (mr_id, complete, checkby) = 
     (SELECT mr_id, complete, checkby 
      FROM ds_moriver
      WHERE mr_fid = ds_fish_check.MR_FID)
      -- and UPLOAD_SESSION_ID = p_mrSessionID )   
 		 WHERE upload_session_id = p_fishSessionID;
 		 
 	END IF; 
      
  END LOOP;
 
 -- insert into ds_fish
      INSERT into ds_fish (mr_id, f_fid, panelhook, bait, species, length, weight,
        last_updated, uploaded_by,fishcount, fin_curl, otolith, rayspine, scale, ftprefix, ftnum, ftmr,
        upload_filename, complete, checkby, upload_session_id)
        SELECT mr_id, f_fid, panelhook, bait, species, length, weight, last_updated, uploaded_by,
            fishcount, fin_curl, otolith, rayspine, scale, ftprefix, ftnum, ftmr, upload_filename,
            complete, checkby, upload_session_id
        FROM ds_fish_check
        WHERE uploaded_by = p_user
         and upload_session_id = p_fishSessionID
         and species IS NOT NULL
         and EXISTS (select mr_fid from ds_moriver
                         where ds_moriver.mr_fid = ds_fish_check.mr_fid);   -- added to only add fish where matching MR record         
         
       
      -- calculate RSD, KN, WR fields after upload
      update ds_fish
        set kn = weight/power(10,-6.2561+3.2932*log(10,length))
        where species = 'PDSG'
        and length IS NOT NULL 
        and weight IS NOT NULL
        and upload_session_id = p_fishSessionID;
        
      update ds_fish
        set Wr = (100*weight/power(10,-6.287+3.330*log(10,length)))
      where species = 'SNSG'
      and nvl(length,0) > 119 and nvl(weight,0) <> 0
      and upload_session_id = p_fishSessionID;
      
      -- NEW - Condition field replacing Kn and Wr - 6/3/19
        update ds_fish
        set condition = weight/power(10,-6.2561+3.2932*log(10,length))
        where species = 'PDSG'
        and nvl(length,0) > 0
        and nvl(weight,0) > 0
        and upload_session_id = p_fishSessionID;
        
      update ds_fish
        set condition = (100*weight/power(10,-6.287+3.330*log(10,length)))
      where species = 'SNSG'
      and nvl(length,0) > 119 and nvl(weight,0) <> 0
      and upload_session_id = p_fishSessionID;
      -----------------------------------------------
      
      update ds_fish
        set RSD =
            (CASE
             WHEN LENGTH <= 199 AND species = 'PDSG' then 'LSS'
             WHEN LENGTH > 199 AND LENGTH <= 329 and species = 'PDSG' THEN 'SS'
             WHEN LENGTH > 329 AND LENGTH <= 629 AND species = 'PDSG' THEN 'SS'
             WHEN LENGTH > 629 AND LENGTH <= 839 AND species = 'PDSG' THEN 'Q'
             WHEN LENGTH > 839 AND LENGTH <= 1039 AND species = 'PDSG' THEN 'P'
             WHEN LENGTH > 1039 AND species = 'PDSG' THEN 'M' 
             WHEN LENGTH <= 149 AND species = 'SNSG' then 'LSS'
             WHEN LENGTH > 149 AND LENGTH <= 249 AND species = 'SNSG' THEN 'SS'
             WHEN LENGTH > 249 AND LENGTH <= 379 AND species = 'SNSG' THEN 'SS'
             WHEN LENGTH > 379 AND LENGTH <= 509 AND species = 'SNSG' THEN 'Q'
             WHEN LENGTH > 509 AND LENGTH <= 639 AND species = 'SNSG' THEN 'P'
             WHEN LENGTH > 640 AND species = 'SNSG' THEN 'M'       
            END)
        WHERE length is not null
        and upload_session_id = p_fishSessionID;        
            
      COMMIT;
 
 ----- Count Fish records uploaded to DS_FISH
 -- count how many new mr records were added to DS_MORIVER
  select count(*)
   into p_fishCntFinal
  from ds_fish
  where upload_session_id = p_fishSessionID; 
  
 -------------------  SUPPLEMENTAL   ---------------------------------------
 
 -- assign corresponding mr_datashee mr_id to supp records mr_id
 update ds_supplemental_check
 set (mr_id,complete,checkby) = 
     (SELECT mr_id, complete, checkby
     FROM ds_moriver
     WHERE ds_supplemental_check.MR_FID = ds_moriver.MR_FID
      and ds_moriver.UPLOAD_SESSION_ID = p_mrSessionID ),
      
     f_id = 
     (SELECT f_id 
     FROM ds_fish
     WHERE ds_supplemental_check.F_FID = ds_fish.F_FID
      and ds_fish.UPLOAD_SESSION_ID = p_fishSessionID )
      
 WHERE upload_session_id = p_suppSessionID;
 
 -- IF adding supplemental that was not in first upload - everything else the same
 FOR Y in (
    select mr_id, complete, checkby
    from ds_supplemental_check
    where upload_session_id = p_suppSessionID)
  
  LOOP
  
    IF y.mr_id IS NULL THEN
    
     update ds_supplemental_check
         set (mr_id,complete,checkby) = 
         (SELECT mr_id, complete, checkby
         FROM ds_moriver
         WHERE ds_supplemental_check.MR_FID = ds_moriver.MR_FID),
          -- and ds_moriver.UPLOAD_SESSION_ID = p_mrSessionID ),
                
         f_id = 
         (SELECT f_id 
         FROM ds_fish
         WHERE ds_supplemental_check.F_FID = ds_fish.F_FID)
          -- and ds_fish.UPLOAD_SESSION_ID = p_fishSessionID )
      
     WHERE upload_session_id = p_suppSessionID;
 		 
 	END IF; 
      
  END LOOP;
 
 -- insert into ds_supplemental
      INSERT into ds_supplemental (mr_id, f_id, f_fid, tagnumber, pitrn, cwtyn, dangler,
        last_updated, uploaded_by, scuteloc, scutenum,
        elhv, elcolor, erhv, ercolor, genetic, genetics_vial_number,
        head, snouttomouth, inter, mouthwidth, m_ib, l_ob, l_ib, r_ib, r_ob, anal, dorsal,
        sex, stage, status, hatchery_origin, recapture, photo,
        genetic_needs, broodstock, hatch_wild, species_id, archive,
        other_tag_info,comments, upload_filename, complete, checkby, upload_session_id,
        scuteloc2, scutenum2)
        SELECT mr_id, f_id, f_fid, tagnumber, pitrn, cwtyn, dangler, last_updated, uploaded_by,
            scuteloc, scutenum, elhv, elcolor, erhv, ercolor, genetic, genetics_vial_number,
            head, snouttomouth, inter, mouthwidth, m_ib, l_ob, l_ib, r_ib, r_ob, anal, dorsal,
            sex, stage, status, hatchery_origin, recapture, photo,
            genetic_needs, broodstock, hatch_wild, species_id, archive,
            other_tag_info, comments, upload_filename, complete, checkby,upload_session_id,
            scuteloc2, scutenum2
        FROM ds_supplemental_check
        WHERE uploaded_by = p_user
         and upload_session_id = p_suppSessionID
         and EXISTS (select f_fid from ds_fish
                         where ds_fish.f_fid = ds_supplemental_check.f_fid);   -- added to only add supplemental where matching MR record ;
         
      COMMIT;
      
  ----- Count Supplemental records uploaded to DS_SUPPLEMENAL
  select count(*)
   into p_suppCntFinal
  from ds_supplemental
  where upload_session_id = p_suppSessionID;
  
  IF p_suppCntFinal > 0 THEN
      
    FOR Z in (
    select s_id, head, inter, m_ib, l_ob, l_ib, r_ib, r_ob, anal, dorsal
    from ds_supplemental
    where upload_session_id = p_suppSessionID
     and uploaded_by = p_user)
  
    LOOP 
    
       IF  (z.head is NOT NULL AND z.head != 0) AND (z.inter is NOT NULL AND z.inter != 0) and (z.m_ib is NOT NULL AND z.m_ib != 0) 
            and (z.l_ob is NOT NULL AND z.l_ob != 0) and
            (z.l_ib is NOT NULL AND z.l_ib != 0) and (z.r_ib is NOT NULL AND z.r_ib != 0) and (z.r_ob is NOT NULL AND z.r_ob != 0) 
            and (z.anal is NOT NULL AND z.anal != 0) and
            (z.dorsal is NOT NULL AND z.dorsal != 0) THEN
 
        update ds_supplemental
         set CI = 6.11 + (0.00000235 *(z.DORSAL)) - (0.177*(z.ANAL)) -
           (0.703*(((z.L_OB + z.R_OB)/2)/ ((z.L_IB + z.R_IB)/2))) - 
           (1.424*(z.HEAD/((z.L_IB + z.R_IB)/2))) + (1.389*(z.HEAD / z.M_IB)) +
           (2.878*(z.INTER/((z.L_IB + z.R_IB)/2))) - (3.258*(z.INTER / z.M_IB))
      
        WHERE s_id = z.s_id;
         
       END IF;
       
       IF  (z.head is NOT NULL AND z.head != 0) AND (z.inter is NOT NULL AND z.inter != 0) and (z.m_ib is NOT NULL AND z.m_ib != 0) 
            and (z.l_ob is NOT NULL AND z.l_ob != 0) and
            (z.l_ib is NOT NULL AND z.l_ib != 0) and (z.r_ib is NOT NULL AND z.r_ib != 0) and (z.r_ob is NOT NULL AND z.r_ob != 0) 
            and (z.anal is NOT NULL AND z.anal != 0) and
            (z.dorsal is NOT NULL AND z.dorsal != 0) THEN
 
        update ds_supplemental
         set 
          MCI = 2.655 - (0.844*(((z.L_OB + z.R_OB)/2)/ ((z.L_IB + z.R_IB)/2))) - 
            (0.749*(z.HEAD/((z.L_IB + z.R_IB)/2))) + (1.292*(z.HEAD/z.M_IB)) +
            (1.874*(z.INTER/((z.L_IB + z.R_IB)/2))) - (3.776*(z.INTER/z.M_IB))
      
        WHERE s_id = z.s_id;
         
       END IF;
         
    END LOOP;
     
     COMMIT;
     
  END IF;
  
  
  EXCEPTION
  
    WHEN noSiteID then
        p_noSiteID_msg := ' '||p_noSite_cnt||' Missouri River datasheets are not associated with an existing or created site.  If new sites were created in the field, please upload PSPA_SITES_DATASHEET...CSV along with the MR, Fish and Supplemental datasheets.';
      
    
      
  

 
 
-- l_session := 0;
-- l_siteMatch := 0;
-- l_siteMatch_tot := 0;
 
 
  END uploadFinal;


END data_upload;