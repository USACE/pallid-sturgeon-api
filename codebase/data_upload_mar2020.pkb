create or replace PACKAGE BODY          data_upload_MAR2020 IS
/*------------------------------------------------------------------------------
Description: This package contains procedures for file upload

Name               Date       Modification
------------------ ---------- ------------
------------------------------------------------------------------------------*/

PROCEDURE p_debug(p_debug_text VARCHAR2, p_apex_session_id NUMBER) IS
BEGIN
/*
INSERT INTO DEBUG_T (ID, debug_text,date_created,apex_sessioN_id)
                                                    VALUES 
                                            (debug_seq.nextval
                                            ,p_debug_text
                                            , SYSDATE  
                                           , p_apex_session_id
                                            );                                 
                                            
                COMMIT;*/
                NULL;
END;

PROCEDURE uploadSearchDatasheet(p_user IN upload_search.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_sfidMatch OUT number,
                                p_fileBrowse IN varchar2 default null,
                                p_upload_session_id IN upload_search.upload_session_id%TYPE
                                ) IS                            
 
v_blob_data       BLOB;
v_blob_len        NUMBER;
v_position        NUMBER;
v_raw_chunk       RAW(10000);
v_char          CHAR(1);
c_chunk_len     number := 1;
v_line          VARCHAR2 (32767) := NULL;
v_data_array      wwv_flow_global.vc_arr2;

v_filename varchar2(200);

-- added for ignoring first row
v_line_count number;
l_last_param_count  NUMBER;

l_session number default 0;

BEGIN

p_cnt := 0;

-- Read data from wwv_flow_files
 IF p_user is not null THEN
 
    select blob_content, filename into v_blob_data, v_filename
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowse;
         
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

       v_line := REPLACE (v_line, ':', '_');
    -- Convert comma to : to use wwv_flow_utilities
       v_line := REPLACE (v_line, ',', ':');
       
       v_line := REPLACE(REPLACE(v_line, CHR(10), ''), CHR(13), '');
       l_last_param_count := INSTR(v_line, ':', 1,14);
       
    -- Convert each column separated by : into array of data
       v_data_array := wwv_flow_utilities.string_to_table (v_line);
    
    -- This is to skip the header
       IF v_line_count > 1 THEN

    -- Insert data into target table
    EXECUTE IMMEDIATE 'INSERT INTO upload_search(SE_FID, DS_ID, SITE_ID, SITE_FID, SEARCH_DATE, RECORDER, SEARCH_TYPE_CODE, START_TIME, START_LATITUDE, START_LONGITUDE, STOP_TIME, 
                          STOP_LATITUDE, STOP_LONGITUDE, TEMP, CONDUCTIVITY, LAST_UPDATED, UPLOAD_SESSION_ID, UPLOADED_BY, UPLOAD_FILENAME)
         VALUES (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19)'
         USING
         v_data_array(1),
         TO_NUMBER(v_data_array(2)),
         TO_NUMBER(v_data_array(3)),
         v_data_array(4),
         TO_DATE(v_data_array(5), 'MM/DD/YYYY'),
         v_data_array(6),
         v_data_array(7),
         REPLACE (v_data_array(8), '_', ':'),--TO_TIMESTAMP(v_data_array(8)), -- Start time, also a date?
         TO_NUMBER(v_data_array(9)),
         TO_NUMBER(v_data_array(10)),
         REPLACE (v_data_array(11), '_', ':'),--TO_TIMESTAMP(v_data_array(11)),
         TO_NUMBER(v_data_array(12)),
         TO_NUMBER(v_data_array(13)),
         TO_NUMBER(v_data_array(14)),
         CASE WHEN LENGTH(v_line) = l_last_param_count THEN NULL -- If no data in the last parameter the array call will fail
         ELSE TO_NUMBER(v_data_array(15))
         END,
         sysdate,
         p_upload_session_id,
         p_user,
         v_filename;
         
         COMMIT;

       END IF; -- v_line_count > 1
         
         -- Clear out
         v_line := NULL;
           
     END IF;
  
  END LOOP;
  
  -- The number of duplicates skipped
  SELECT count(*) INTO p_sfidMatch FROM upload_search us
   WHERE us.upload_session_id = p_upload_session_id AND (se_fid, site_id) IN (SELECT se_fid, site_id FROM ds_search);
   
  -- The number of new rows added
  SELECT COUNT(*) INTO p_cnt
    FROM upload_search us
   WHERE us.upload_session_id = p_upload_session_id AND (se_fid, site_id) NOT IN (SELECT se_fid, site_id FROM ds_search);
   
 DELETE FROM APEX_APPLICATION_FILES 
  WHERE name = p_fileBrowse;
  
 COMMIT;
END uploadSearchDatasheet;


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
 p_debug('should copy to ds_supplemental around here. l_session = ' || l_session, NULL);
 
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

PROCEDURE uploadTelemetryDatasheet (
                                p_user IN upload_telemetry_fish.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_tfidMatch OUT number,
                                p_fileBrowse IN varchar2 default null,
                                p_upload_session_id IN upload_telemetry_fish.upload_session_id%TYPE
                            ) IS

v_blob_data       BLOB;
v_blob_len        NUMBER;
v_position        NUMBER;
v_raw_chunk       RAW(10000);
v_char          CHAR(1);
c_chunk_len     number := 1;
v_line          VARCHAR2 (32767) := NULL;
v_data_array      wwv_flow_global.vc_arr2;

v_filename varchar2(200);

-- added for ignoring first row
v_line_count number;

BEGIN

-- Read data from wwv_flow_files
 IF p_user is not null THEN

    select blob_content into v_blob_data
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowse;
         
    select filename into v_filename
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowse;
           
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

       v_line := REPLACE (v_line, ':', '_');
    -- Convert comma to : to use wwv_flow_utilities
       v_line := REPLACE (v_line, ',', ':');
       v_line := REPLACE(REPLACE(v_line, CHR(10), ''), CHR(13), '');
       
    -- Convert each column separated by : into array of data
       v_data_array := wwv_flow_utilities.string_to_table (v_line);

    -- This is to skip the header
       IF v_line_count > 1 THEN

    -- Insert data into target table
       EXECUTE IMMEDIATE 'INSERT INTO upload_telemetry_fish(T_FID, SE_FID, BEND, RADIO_TAG_NUM, FREQUENCY_ID_CODE, CAPTURE_TIME, CAPTURE_LATITUDE, CAPTURE_LONGITUDE,
         POSITION_CONFIDENCE, MACRO_ID, MESO_ID, DEPTH, TEMP, CONDUCTIVITY, TURBIDITY, SILT, SAND, GRAVEL, COMMENTS, LAST_UPDATED, UPLOAD_SESSION_ID, UPLOADED_BY, UPLOAD_FILENAME)
         VALUES (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23)'
         USING 
         v_data_array(1),
         v_data_array(2),
         TO_NUMBER(v_data_array(3)),
         TO_NUMBER(v_data_array(4)),
         TO_NUMBER(v_data_array(5)),
         REPLACE (v_data_array(6), '_', ':'), -- This is capture_time
         TO_NUMBER(v_data_array(7)),
         TO_NUMBER(v_data_array(8)),         
         TO_NUMBER(v_data_array(9)),
         v_data_array(10), -- Macro_ID
         v_data_array(11),
         TO_NUMBER(v_data_array(12)),
         TO_NUMBER(v_data_array(13)),
         TO_NUMBER(v_data_array(14)),
         TO_NUMBER(v_data_array(15)),
         TO_NUMBER(v_data_array(16)),         
         TO_NUMBER(v_data_array(17)),  
         TO_NUMBER(v_data_array(18)),  
         v_data_array(19),
         sysdate,
         p_upload_session_id,
         p_user,
         v_filename;
         
         COMMIT;

       END IF; -- v_line_count > 1
         
        -- Clear out
           v_line := NULL;
           
     END IF;
  
  END LOOP;
  
  -- The number of duplicates skipped
  SELECT count(*) INTO p_tfidMatch FROM upload_telemetry_fish utf
   WHERE utf.upload_session_id = p_upload_session_id AND (t_fid, se_fid) IN (SELECT t_fid, se_fid FROM ds_telemetry_fish);
    
  -- The number of new rows added
  SELECT count(*) INTO p_cnt FROM upload_telemetry_fish utf
   WHERE utf.upload_session_id = p_upload_session_id AND (t_fid, se_fid) NOT IN (SELECT t_fid, se_fid FROM ds_telemetry_fish);
   
 DELETE from APEX_APPLICATION_FILES 
  WHERE name = p_fileBrowse;
  
 COMMIT; 
END uploadTelemetryDatasheet;

PROCEDURE uploadSiteDatasheetCheck (
                                p_user IN upload_sites.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_siteMatch OUT number,
                                p_fileBrowseSite IN varchar2 default null,
                                p_upload_session_id IN upload_sites.upload_session_id%TYPE
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

l_fish_max_session number default 0;
l_mr_max_session number default 0;
l_supp_max_session number default 0;

l_siteMatch number default 0;
l_siteMatch_tot number default 0;


BEGIN

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
   
   
   BEGIN
   
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
         p_upload_session_id,
         p_user,
         v_filename;
         
         /*  INSERT INTO DEBUG_T (ID, debug_text,date_created)
           VALUES (debug_seq.nextval
                                ,  '1= Number' || TO_CHAR(v_data_array(1))
                                                 
                , SYSDATE   );
                COMMIT;
         */
         /*
         EXCEPTION
          WHEN others tHEN 
          INSERT INTO DEBUG_T (ID, debug_text,date_created)
           VALUES (debug_seq.nextval
                                ,  '1= Number' || TO_CHAR(v_data_array(1))
                                                 
                , SYSDATE   );
                COMMIT;
              */
         END;
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
    WHERE upload_session_id = p_upload_session_id
     AND uploaded_by = p_user
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
         and upload_session_id = p_upload_session_id
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
 
 l_siteMatch := 0;
 l_siteMatch_tot := 0;
 
 
END uploadSiteDatasheetCheck;

------------

PROCEDURE uploadMRdatasheetCheck (
                                p_user IN upload_mr.uploaded_by%TYPE,
                                -- p_complete IN upload_mr.complete%TYPE,
                                p_checkby IN upload_mr.checkby%TYPE,
                                p_cnt OUT number,
                                p_mrfidMatch OUT number,
                                p_fileBrowseMR IN varchar2 default null,  
                                p_upload_session_id IN upload_mr.upload_session_id%TYPE
                            ) IS

t_d_api VARCHAR2(40) DEFAULT 'UploadMRdatasheedCheck';

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

l_max_mr_id number default 0;
-- l_site_id number default 0;
l_mr_max_session number default 0;
l_site_max_session number default 0;

l_mrDuplicate number default 0;
l_mrDuplicateTot number default 0;


BEGIN

p_debug ( t_d_api ||  ' Inside uploadMRdatasheetCheck start.' || ' CSV line' || v_line_count, NULL);

p_cnt := 0;
p_mrfidMatch := 0;

-- test

-- Read data from wwv_flow_files
 IF p_user is not null THEN
         
         p_debug ( t_d_api ||     ' Inside uploadMRdatasheetCheck getting the wwv_Flow_Files content p_User = ' || p_user || ' file = ' || p_fileBrowseMR   , NULL);

    SELECT blob_content
                 , filename
     INTO v_blob_data
             , v_filename
    FROM APEX_APPLICATION_FILES 
    WHERE name = p_fileBrowseMR;
        

 END IF;
 
 p_debug(  t_d_api ||   ' Past p_user_check.' || 'p_user =' || p_user, NULL);
 
v_blob_len := dbms_lob.getlength(v_blob_data);
v_position := 1;
v_line_count := 0;

 p_debug(  t_d_api ||   ' Going into reading the CSV. v_line_count = ' || v_line_count
                        || ' v_blob_len = ' || v_blob_len
                        || ' v_position = ' || v_blob_len                        
                        , NULL );
                
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

        p_debug( 'After line_count > 1 = ' || v_line_count
                 || ' about to start inserting into the upload_mr table' , NULL );

    -- Insert data into target table

insert into upload_mr (site_id, site_fid, mr_fid, se_Field_id, season, setdate, subsample, subsamplepass, 
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
                                             values
( to_number(v_data_array(1)), --site_id, XLS col1
         v_data_array(2), --site_fid, XLS col2
         v_data_array(3), --mr_Fid, XLS col3
         v_data_array(4), --se_fileD_id col4
         v_data_array(5), --season  xls, col5
         to_date(v_data_array(6),'mm/dd/yyyy'),
         to_number(v_data_array(7)),
         to_number(v_data_array(8)),
         v_data_array(9),
         v_data_array(10),
         
         v_data_array(11),
         v_data_array(12),
         to_number(v_data_array(13)),
         to_number(v_data_array(14)),
         to_number(v_data_array(15)),
         to_number(v_data_array(16)),
         
         to_number(v_data_array(17)),
         to_number(v_data_array(18)),
         
         to_number(v_data_array(19)),
         v_data_array(20),
         
         v_data_array(21),
         to_number(v_data_array(22)),
         to_number(v_data_array(23)),
         
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
         
         v_data_array(41),
         to_number(v_data_array(42)),
         to_number(v_data_array(43)),
         v_data_array(44),
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
         to_number(v_data_array(64)),         
         v_data_array(67),  --  
         sysdate,
         p_upload_session_id,
         p_user,
         v_filename,
         p_checkby,
         v_data_array(65),
         v_data_array(66)
);
       
         p_debug ('Past insert... a new row should be in the table. ' || ' CSV line' || v_line_count, NULL );
         
    /* original John Code -- commented out 5APR2020 by JDK
       EXECUTE IMMEDIATE 'insert into upload_mr (site_id, site_fid, mr_fid, season, setdate
       , subsample, subsamplepass,  subsamplen, recorder, gear
       , gear_type, temp, turbidity, conductivity, do,        distance
       , width, netrivermile, structurenumber, usgs, riverstage
       , discharge, u1, u2, u3, u4
       , u5, u6, u7, macro, meso
       , habitatrn, qc, micro_structure, structure_flow, structure_mod
       , set_site_1, set_site_2, set_site_3, starttime, startlatitude
       , startlongitude, stoptime, stoplatitude, stoplongitude,
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
    */

       END IF; -- v_line_count > 1

              p_debug (  'About to Clear the line = ' || v_line_count || '. ' , NULL);
         
        -- Clear out
           v_line := NULL;
           
     END IF;
  END LOOP;

COMMIT;

p_debug(  'Past parsing CSV = ' || v_line_count || '. ' , NULL);

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
     and upload_session_id = p_upload_session_id)
  
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
            upload_session_id, site_fid, se_field_id
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
            upload_session_id, site_fid, se_field_id
        FROM UPLOAD_MR
        WHERE uploaded_by = p_user
         and upload_session_id = p_upload_session_id
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
 
END uploadMRdatasheetCheck;

-----

PROCEDURE uploadFishDatasheetCheck (
                                p_user IN upload_fish.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_ffidMatch OUT number,
                                p_fileBrowseF IN varchar2 default null,
                                p_upload_session_id IN upload_fish.upload_session_id%TYPE
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

 l_fish_max_session number default 0;
 l_mr_max_session number default 0;
 
l_fDuplicate number default 0;
l_fDuplicateTot number default 0;


BEGIN

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
         p_upload_session_id,
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
     and upload_session_id = p_upload_session_id)
  
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
         and upload_session_id = p_upload_session_id
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
  
END uploadFishDatasheetCheck;

-------
PROCEDURE uploadProcedureDatasheet (            
                                p_user IN upload_procedure.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_pfidMatch OUT number,
                                p_fileBrowseP IN varchar2 default null, 
                                p_upload_session_id IN upload_procedure.upload_session_id%TYPE
                          ) IS

v_blob_data     BLOB;
v_blob_len      NUMBER;
v_position      NUMBER;
v_raw_chunk     RAW(10000);
v_char          CHAR(1);
c_chunk_len     number := 1;
v_line          VARCHAR2 (32767) := NULL;
v_data_array    wwv_flow_global.vc_arr2;

v_filename varchar2(200);

-- added for ignoring first row
v_line_count number;

l_p_id  number;
 
BEGIN

p_cnt := 0;
p_pfidMatch := 0;

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
         
    select blob_content, filename into v_blob_data, v_filename
    from APEX_APPLICATION_FILES 
    where name = p_fileBrowseP;
               
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

       v_line := REPLACE (v_line, ':', '_');
    -- Convert comma to : to use wwv_flow_utilities
       v_line := REPLACE (v_line, ',', ':');

    -- Convert each column separated by : into array of data
       v_data_array := wwv_flow_utilities.string_to_table (v_line);

       IF v_line_count > 1 THEN
       /*
       EXECUTE IMMEDIATE 'insert into upload_procedure (P_FID, S_FID, S_ID, PURPOSE_CODE, PROCEDURE_DATE, PROCEDURE_START_TIME, PROCEDURE_END_TIME, PROCEDURE_BY, 
                                                        ANTIBIOTIC_INJECTION_IND, PHOTO_DORSAL_IND, PHOTO_VENTRAL_IND, PHOTO_LEFT_IND)
                                                        /*OLD_RADIO_TAG_NUM, OLD_FREQUENCY_ID, DST_SERIAL_NUM, DST_START_DATE, DST_START_TIME, DST_REIMPLANT_IND, NEW_RADIO_TAG_NUM,
                                                        NEW_FREQUENCY_ID, SEX_CODE, BLOOD_SAMPLE_IND, EGG_SAMPLE_IND, COMMENTS, FISH_HEALTH_COMMENTS,
                                                        EVAL_LOCATION_CODE, SPAWN_CODE, VISUAL_REPRO_STATUS_CODE, ULTRASOUND_REPRO_STATUS_CODE,
                                                        EXPECTED_SPAWN_YEAR, ULTRASOUND_GONAD_LENGTH, GONAD_CONDITION,
                                                        LAST_UPDATED, UPLOAD_SESSION_ID, UPLOADED_BY, UPLOAD_FILENAME )                                                        
         values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10)' --,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23,:24,:25,:26,:27,:28,:29,:30,:31)' --,:31,:32,:33,:34)'
         -- Line breaks below correspond to line breaks above
         USING 
        -- null, -- l_p_id,
         v_data_array(1),
         v_data_array(2),
         p_suppSessionID, -- this is the system id for supplemental, populated later
         v_data_array(3),         
         v_data_array(4),
         v_data_array(5),
         v_data_array(6),
         v_data_array(7),
               
         to_number(v_data_array(8)),         
         to_number(v_data_array(9)),     
         to_number(v_data_array(10)),
         to_number(v_data_array(11));
         /*
         to_number(v_data_array(11)),
         to_number(v_data_array(12)),
         to_number(v_data_array(13)),
         v_data_array(14),
         v_data_array(15),
         to_number(v_data_array(16)),         
         to_number(v_data_array(17)),
         
         to_number(v_data_array(18)),      
         v_data_array(19),        -- T sex 
         to_number(v_data_array(20)),         
         to_number(v_data_array(21)),
         v_data_array(22),
         v_data_array(23),         
               
         v_data_array(24),
         v_data_array(25),
         v_data_array(26),
         v_data_array(27),
                  
         to_number(v_data_array(28)),
         to_number(v_data_array(29)),
         v_data_array(30),
         
         sysdate,
         p_pSessionID,
         p_user,
         v_filename;
         */
         commit;
         /*
         track_ben(55, 'upload procedure - v_line', v_line_count || ' - ' || v_line);
         track_ben(55, 'upload procedure - 21', v_line_count || ' - ' || v_data_array(21));
         track_ben(55, 'upload procedure - 22', v_line_count || ' - ' || v_data_array(22));
         track_ben(55, 'upload procedure - 29', v_line_count || ' - ' || v_data_array(29));
         track_ben(55, 'upload procedure - 30', v_line_count || ' - ' || v_data_array(30));
         */
         
    -- Insert data into target table
       EXECUTE IMMEDIATE 'insert into upload_procedure (ID, F_FID, PURPOSE_CODE, PROCEDURE_DATE, PROCEDURE_START_TIME, PROCEDURE_END_TIME, PROCEDURE_BY, 
                                                        ANTIBIOTIC_INJECTION_IND, PHOTO_DORSAL_IND, PHOTO_VENTRAL_IND, PHOTO_LEFT_IND,
                                                        OLD_RADIO_TAG_NUM, OLD_FREQUENCY_ID, DST_SERIAL_NUM, DST_START_DATE, DST_START_TIME, DST_REIMPLANT_IND, NEW_RADIO_TAG_NUM,
                                                        NEW_FREQUENCY_ID, SEX_CODE, BLOOD_SAMPLE_IND, EGG_SAMPLE_IND, COMMENTS, FISH_HEALTH_COMMENTS,
                                                        EVAL_LOCATION_CODE, SPAWN_CODE, VISUAL_REPRO_STATUS_CODE, ULTRASOUND_REPRO_STATUS_CODE,
                                                        EXPECTED_SPAWN_YEAR, ULTRASOUND_GONAD_LENGTH, GONAD_CONDITION,
                                                        LAST_UPDATED, UPLOAD_SESSION_ID, UPLOADED_BY, UPLOAD_FILENAME )                                                        
         values (:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17,:18,:19,:20,:21,:22,:23,:24,:25,:26,:27,:28,:29,:30,:31,:32,:33,:34,:35)'
         -- Line breaks below correspond to line breaks above
         USING 
         procedure_seq.nextval, 
         v_data_array(1),
         --v_data_array(2),
         --p_suppSessionID, -- this is the system id for supplemental, populated later
         v_data_array(3),    
         v_data_array(4),  -- This should be converted TO_DATE
         REPLACE (v_data_array(5), '_', ':'),
         REPLACE (v_data_array(6), '_', ':'),
         v_data_array(7),
                
         to_number(v_data_array(8)),         
         to_number(v_data_array(9)),     
         to_number(v_data_array(10)),
         to_number(v_data_array(11)),  
         
         to_number(v_data_array(12)),
         to_number(v_data_array(13)),
         to_number(v_data_array(14)),         
         v_data_array(15),
         REPLACE (v_data_array(16), '_', ':'),
         to_number(v_data_array(17)),         
         to_number(v_data_array(18)),
         
         to_number(v_data_array(19)),      
         v_data_array(20),        -- T sex 
         to_number(v_data_array(21)),         
         to_number(v_data_array(22)),
         v_data_array(23),
         v_data_array(24),         
               
         v_data_array(25),
         v_data_array(26),
         v_data_array(27),
         v_data_array(28),
                  
         to_number(v_data_array(29)),
         to_number(v_data_array(30)),
         v_data_array(31),
         
         sysdate,
         p_upload_session_id,
         p_user,
         v_filename;
         
         commit;

       END IF; -- v_line_count > 1
         
        -- Clear out
           v_line := NULL;
           
     END IF;
  
  END LOOP;

  v_line_count := 0; 
    
SELECT COUNT(*) INTO p_pfidMatch
  FROM ds_procedure dsp
 INNER JOIN upload_procedure up ON (dsp.f_fid = up.f_fid);
   
   -- count how many new mr records were added to DS_MORIVER
  select count(*)
   into p_cnt
   from upload_procedure
  where upload_session_id = p_upload_session_id; 
    
DELETE from APEX_APPLICATION_FILES 
 where name = p_fileBrowseP;  
 
/*
  -- If there are p_fid that already exist in DS_PROCEDURE then don't upload any records to DS_PROCEDURE_CHECK
  IF p_pfidMatch = 0 THEN
   
      -- insert into ds_procedure
      INSERT INTO ds_procedure (ID, P_FID, S_FID, S_ID, PURPOSE_CODE, PROCEDURE_DATE, PROCEDURE_START_TIME, PROCEDURE_END_TIME, PROCEDURE_BY, 
                                                        ANTIBIOTIC_INJECTION_IND, PHOTO_DORSAL_IND, PHOTO_VENTRAL_IND, PHOTO_LEFT_IND,
                                                        OLD_RADIO_TAG_NUM, OLD_FREQUENCY_ID, DST_SERIAL_NUM, DST_START_DATE, DST_START_TIME, DST_REIMPLANT_IND, NEW_RADIO_TAG_NUM,
                                                        NEW_FREQUENCY_ID, SEX_CODE, BLOOD_SAMPLE_IND, EGG_SAMPLE_IND, COMMENTS, FISH_HEALTH_COMMENTS,
                                                        EVAL_LOCATION_CODE, SPAWN_CODE, VISUAL_REPRO_STATUS_CODE, ULTRASOUND_REPRO_STATUS_CODE,
                                                        EXPECTED_SPAWN_YEAR, ULTRASOUND_GONAD_LENGTH, GONAD_CONDITION,
                                                        LAST_UPDATED, UPLOAD_SESSION_ID, UPLOADED_BY, UPLOAD_FILENAME )
        SELECT ID, P_FID, S_FID, S_ID, PURPOSE_CODE, PROCEDURE_DATE, PROCEDURE_START_TIME, PROCEDURE_END_TIME, PROCEDURE_BY, 
                                                        ANTIBIOTIC_INJECTION_IND, PHOTO_DORSAL_IND, PHOTO_VENTRAL_IND, PHOTO_LEFT_IND,
                                                        OLD_RADIO_TAG_NUM, OLD_FREQUENCY_ID, DST_SERIAL_NUM, DST_START_DATE, DST_START_TIME, DST_REIMPLANT_IND, NEW_RADIO_TAG_NUM,
                                                        NEW_FREQUENCY_ID, SEX_CODE, BLOOD_SAMPLE_IND, EGG_SAMPLE_IND, COMMENTS, FISH_HEALTH_COMMENTS,
                                                        EVAL_LOCATION_CODE, SPAWN_CODE, VISUAL_REPRO_STATUS_CODE, ULTRASOUND_REPRO_STATUS_CODE,
                                                        EXPECTED_SPAWN_YEAR, ULTRASOUND_GONAD_LENGTH, GONAD_CONDITION,
                                                        LAST_UPDATED, UPLOAD_SESSION_ID, UPLOADED_BY, UPLOAD_FILENAME 
        FROM UPLOAD_PROCEDURE
        WHERE upload_session_id = l_p_session_id;   
         
      COMMIT;
      
  END IF;
  
  -- count how many new mr records were added to DS_MORIVER
  select count(*)
   into p_cnt
   from ds_procedure
  where upload_session_id = l_p_session_id; 
    
DELETE from APEX_APPLICATION_FILES 
 where name = p_fileBrowseP;
 */
END uploadProcedureDatasheet;

-----------
PROCEDURE uploadSuppDatasheetCheck (
                                p_user IN upload_supplemental.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_sfidMatch OUT number,
                                p_fileBrowseSup IN varchar2 default null,
                                p_upload_session_id IN upload_supplemental.upload_session_id%TYPE,
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
         p_upload_session_id,
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
     and upload_session_id = p_upload_session_id)
  
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
         and upload_session_id = p_upload_session_id
         and NOT EXISTS (select f_fid from ds_supplemental
                         where ds_supplemental.f_fid = upload_supplemental.f_fid); -- added to keep duplicates from being entered
         
      COMMIT;
      
      
  -- NEW SECTION -- To catch '+' in Tag number for Supplemental - 1/6/14
  -- Loop through upload_supplemental records for this session
  FOR Z in (
    SELECT f_fid, tagnumber
    FROM ds_supplemental_check
    WHERE uploaded_by = p_user
     and upload_session_id = p_upload_session_id)
  
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
  
END uploadSuppDatasheetCheck;


PROCEDURE uploadSearchDatasheetCheck(p_user IN upload_supplemental.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_searchMatch OUT number,
                                p_fileBrowse IN varchar2 default null,
                                p_upload_session_id IN upload_search.upload_session_id%TYPE
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

l_max_mr_id number default 0;
-- l_site_id number default 0;
l_mr_max_session number default 0;
l_site_max_session number default 0;

l_mrDuplicate number default 0;
l_mrDuplicateTot number default 0;

BEGIN
NULL;
--set variables for this procedure
    p_cnt                       := 0;
    p_searchMatch   := 0;

-- Read data from wwv_flow_files
 IF p_user is not null THEN
    SELECT blob_content
                 , filename
     INTO v_blob_data
              , v_filename
    FROM apex_application_files 
    WHERE name = p_fileBrowse;
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
      
   BEGIN
   INSERT INTO ds_Search (se_fid,ds_id, site_Id,site_fid,search_date, recorder
   , search_type_code
   , start_time,start_latitude  , start_longitude
   , stop_time,stop_latitude    ,  stop_longitude
   ,temp, conductivity
   ,se_id
   ) VALUES
   ( v_data_array(1)
   ,TO_NUMBER(v_data_array(2))
   ,TO_NUMBER(v_data_array(3))
   ,v_data_array(4)
   ,TO_DATE(v_data_array(5),'mm/dd/yyyy')
    ,v_data_array(6)
    ,v_data_array(7)
    ,TO_DATE(v_data_array(5) || ' ' || v_data_array(8) ,'mm/dd/yyyy hh24:mi:ss')
    ,TO_NUMBER(v_data_array(9))
    ,TO_NUMBER(v_data_array(10))
    ,TO_DATE(v_data_array(5) || ' ' || v_data_array(11) ,'mm/dd/yyyy hh24:mi:ss')
    ,TO_NUMBER(v_data_array(12))
    ,TO_NUMBER(v_data_array(13)) 
    ,TO_NUMBER(v_data_array(14))
    ,TO_NUMBER(v_data_array(15))
   , SEARCH_SEQ.NEXTVAL --se_id  PK
   );
   
                COMMIT;
    
         END;
      END IF; -- v_line_count > 1
 
-- Clear out
   v_line := NULL;
   
 END IF;
  
 END LOOP;

  v_line_count := 0;

END uploadSearchDatasheetCheck;

PROCEDURE uploadFinal (
                                p_user IN ds_sites_check.uploaded_by%TYPE,
                                p_site_cnt_final OUT number,
                                p_mr_cnt_final OUT number,
                                p_fishCntFinal OUT number,
                                p_searchCntFinal OUT NUMBER,
                                p_suppCntFinal OUT number,
                                p_telemetryCntFinal OUT NUMBER,
                                p_procedureCntFinal OUT NUMBER,
                                p_plus_cnt in number,
                                p_noSite_cnt OUT number,
                                p_noSiteID_msg OUT varchar2,
                                p_upload_session_id IN upload_search.upload_session_id%TYPE
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
  p_debug ('uploadFinal' || ' starting. ', NULL);
      -- insert into ds_sites - only those with site_FID
      INSERT into ds_sites (site_fid, year, fieldoffice, project_id,
            segment_id, season, bend, bendrn,
            last_updated, uploaded_by,upload_filename,upload_session_id, sample_unit_type)
        
        SELECT site_fid, year, fieldoffice, project_id,
            segment_id, season, bend, bendrn,
            last_updated, uploaded_by,upload_filename, upload_session_id, 'B'
        FROM ds_sites_check
        WHERE uploaded_by = p_user
         and upload_session_id = p_upload_session_id
         and site_fid IS NOT NULL;
      
      p_debug ('uploadFinal' || ' Past inserting  ds_sites. ', NULL);   
  
  -- count how many new sites were added to DS_SITES
  select count(*)
   into p_site_cnt_final
  from ds_sites
  where upload_session_id = p_upload_session_id;
  
        p_debug ('uploadFinal' || ' Past count of DS_Sites. ', NULL);

  -- Missouri River Data Sheets  
                       
    -- if uploading a MR data sheet associated with new site
    FOR x IN
        (SELECT SITE_ID
         FROM DS_MORIVER_CHECK
         WHERE upload_session_id = p_upload_session_id
         )
     LOOP
     
             p_debug ('inside Select Site_id From DS_MORIVER_CHECK Loop' || '', NULL);
        
        IF x.SITE_ID IS NULL OR x.SITE_ID = 0 THEN
        
             update DS_MORIVER_CHECK
             set site_id = 
                 (SELECT site_id
                  FROM DS_SITES
                  WHERE site_fid = ds_moriver_check.site_fid
                    and upload_session_id = p_upload_session_id
                  )
             WHERE upload_session_id = p_upload_session_id
                  AND site_id = x.site_id;
        END IF;
     
     END LOOP; 
     
     -- count how many records in DS_MORIVER_CHECK don't have SITE_ID after process above.
     -- If some records don't have SITE_ID then they tried to upload MR data with field-added site without uploading new sites
     SELECT count(*)
        INTO p_noSite_cnt     
     FROM  DS_MORIVER_CHECK
     WHERE SITE_ID IS NULL
     and upload_session_id = p_upload_session_id;
     
                  p_debug ('Got count of no site count p_noSite_cnd = ' || p_noSite_cnt, NULL);

     
     
     -- Exception - if p_noSite_cnt is > 0 then do exception and jump to bottom
     IF NVL(p_noSite_cnt,0) > 0 THEN
        NULL;
        --raise noSiteID;
     END IF;

     
     -- update DS_MORIVER_CHECK with COMPLETE field per MR_FID based on count of BAFI per MR_FID
     -- 11/12/14
     FOR X IN
        (SELECT MR_FID
         FROM DS_MORIVER_CHECK
         WHERE upload_session_id = p_upload_session_id)  
     LOOP
     
        FOR Y IN
            (SELECT count(species) as speciesCount
             FROM DS_FISH_CHECK
             WHERE upload_session_id = p_upload_session_id
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
                WHERE upload_session_id = p_upload_session_id
                 and mr_fid = x.mr_fid;            
            
            END IF;                  
        
        END LOOP;     
     
     END LOOP;     
     
  p_debug ('processFinalUpload Inserting Records into the main tables', NULL);

  -- Insert any new rows into DS_SEARCH
  INSERT INTO DS_SEARCH (SE_ID, SEARCH_DATE, RECORDER, SEARCH_TYPE_CODE, START_TIME, START_LATITUDE, START_LONGITUDE, STOP_TIME, STOP_LATITUDE, STOP_LONGITUDE, SE_FID,
                         DS_ID, SITE_ID, SITE_FID, TEMP, CONDUCTIVITY, LAST_UPDATED, UPLOAD_SESSION_ID, UPLOADED_BY, UPLOAD_FILENAME)  
  SELECT search_seq.nextval, SEARCH_DATE, RECORDER, SEARCH_TYPE_CODE, START_TIME, START_LATITUDE, START_LONGITUDE, STOP_TIME, STOP_LATITUDE, STOP_LONGITUDE, SE_FID,
                         DS_ID, SITE_ID, SITE_FID, TEMP, CONDUCTIVITY, LAST_UPDATED, UPLOAD_SESSION_ID, UPLOADED_BY, UPLOAD_FILENAME
    FROM upload_search us
   WHERE us.upload_session_id = p_upload_session_id AND (se_fid, site_id) NOT IN (SELECT se_fid, site_id FROM ds_search);
  
  -- count how many new search effort records were added to DS_SEARCH
  SELECT COUNT(*)
    INTO p_searchCntFinal
    FROM ds_search
   WHERE upload_session_id = p_upload_session_id;

  -- Insert any new rows into DS_TELEMETRY_FISH
  INSERT INTO ds_telemetry_fish(t_id, T_FID, SE_ID, BEND, RADIO_TAG_NUM, FREQUENCY_ID_CODE, CAPTURE_TIME, CAPTURE_LATITUDE, CAPTURE_LONGITUDE,
         POSITION_CONFIDENCE, MACRO_ID, MESO_ID, DEPTH, TEMP, CONDUCTIVITY, TURBIDITY, SILT, SAND, GRAVEL, COMMENTS, LAST_UPDATED, UPLOAD_SESSION_ID, UPLOADED_BY, UPLOAD_FILENAME) 
  SELECT telemetry_id_seq.nextval, T_FID, s.se_id, BEND, RADIO_TAG_NUM, FREQUENCY_ID_CODE, CAPTURE_TIME, CAPTURE_LATITUDE, CAPTURE_LONGITUDE,
         POSITION_CONFIDENCE, MACRO_ID, MESO_ID, DEPTH, utf.TEMP, utf.CONDUCTIVITY, TURBIDITY, SILT, SAND, GRAVEL, COMMENTS, utf.LAST_UPDATED, utf.UPLOAD_SESSION_ID, utf.UPLOADED_BY, utf.UPLOAD_FILENAME
    FROM upload_telemetry_fish utf
   INNER JOIN ds_search s ON (utf.se_fid = s.se_fid AND s.upload_session_id = p_upload_session_id)
   WHERE utf.upload_session_id = p_upload_session_id AND (utf.t_fid, s.se_id) NOT IN (SELECT tf.t_fid, tf.se_id FROM ds_telemetry_fish tf);

  -- count how many new telemetry records were added to DS_TELEMETRY_FISH
  SELECT COUNT(*)
    INTO p_telemetryCntFinal
    FROM ds_telemetry_fish
   WHERE upload_session_id = p_upload_session_id;
     
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
            upload_session_id,
            se_id
            )
        SELECT mr.site_id, mr_fid, season, setdate, subsample, subsamplepass, mr.comments, mr.last_updated, mr.uploaded_by,
            subsamplen, mr.recorder, gear, gear_type, mr.temp, turbidity, mr.conductivity, do,
            distance, width, netrivermile, structurenumber, usgs, riverstage, discharge,
            u1, u2, u3, u4, u5, u6, u7, macro, meso, habitatrn, qc,
            micro_structure, structure_flow, structure_mod, set_site_1, set_site_2, set_site_3,
            micro, 
            starttime, startlatitude, startlongitude, 
            stoptime, stoplatitude, stoplongitude,
            depth1, velocitybot1, velocity08_1, velocity02or06_1,
            depth2, velocitybot2, velocity08_2, velocity02or06_2,
            depth3, velocitybot3, velocity08_3, velocity02or06_3,
            watervel, cobble, organic, silt, sand, gravel, mr.upload_filename,
            complete, checkby, 
            no_turbidity, no_velocity, 
            mr.upload_session_id, s.se_id
        FROM DS_MORIVER_CHECK mr
        LEFT OUTER JOIN ds_search s ON (mr.se_field_id = s.se_fid AND s.upload_session_id = p_upload_session_id)
        WHERE mr.uploaded_by = p_user
         and mr.upload_session_id = p_upload_session_id;
          
     END IF;
 
  p_debug ('Around line 2530', NULL);

  -- count how many new mr records were added to DS_MORIVER
  select count(*)
   into p_mr_cnt_final
  from ds_moriver
  where upload_session_id = p_upload_session_id;
  
  ----------------------- FISH  -------------------
 update ds_fish_check
 set (mr_id, complete, checkby) = 
     (SELECT mr_id, complete, checkby 
     FROM ds_moriver
     WHERE mr_fid = ds_fish_check.MR_FID
      and UPLOAD_SESSION_ID = p_upload_session_id )   
 WHERE upload_session_id = p_upload_session_id;
 p_debug ('Around line 2546. p_fishSessionID =' || p_upload_session_id, NULL);
 -- If adding fish that weren't in original upload...
 
 FOR Y in (
    select mr_id, complete, checkby
    from ds_fish_check
    where upload_session_id = p_upload_session_id)
  
  LOOP
  
    IF y.mr_id IS NULL THEN
    
     update ds_fish_check
 			set (mr_id, complete, checkby) = 
     (SELECT mr_id, complete, checkby 
      FROM ds_moriver
      WHERE mr_fid = ds_fish_check.MR_FID)
      -- and UPLOAD_SESSION_ID = p_mrSessionID )   
 		 WHERE upload_session_id = p_upload_session_id;
 		 
 	END IF; 
      
  END LOOP;
  
  p_debug ('Around line 2568', NULL);
 -- insert into ds_fish
      INSERT into ds_fish (mr_id, f_fid, panelhook, bait, species, length, weight,
        last_updated, uploaded_by,fishcount, fin_curl, otolith, rayspine, scale, ftprefix, ftnum, ftmr,
        upload_filename, complete, checkby, upload_session_id)
        SELECT mr_id, f_fid, panelhook, bait, species, length, weight, last_updated, uploaded_by,
            fishcount, fin_curl, otolith, rayspine, scale, ftprefix, ftnum, ftmr, upload_filename,
            complete, checkby, upload_session_id
        FROM ds_fish_check
        WHERE uploaded_by = p_user
         and upload_session_id = p_upload_session_id
         and species IS NOT NULL
         and EXISTS (select mr_fid from ds_moriver
                         where ds_moriver.mr_fid = ds_fish_check.mr_fid);   -- added to only add fish where matching MR record         

  p_debug ('Around line 2582', NULL);

      -- calculate RSD, KN, WR fields after upload
      update ds_fish
        set kn = weight/power(10,-6.2561+3.2932*log(10,length))
        where species = 'PDSG'
        and length IS NOT NULL 
        and weight IS NOT NULL
        and upload_session_id = p_upload_session_id;
        
      update ds_fish
        set Wr = (100*weight/power(10,-6.287+3.330*log(10,length)))
      where species = 'SNSG'
      and nvl(length,0) > 119 and nvl(weight,0) <> 0
      and upload_session_id = p_upload_session_id;
      
      -- NEW - Condition field replacing Kn and Wr - 6/3/19
        update ds_fish
        set condition = weight/power(10,-6.2561+3.2932*log(10,length))
        where species = 'PDSG'
        and nvl(length,0) > 0
        and nvl(weight,0) > 0
        and upload_session_id = p_upload_session_id;
        
      update ds_fish
        set condition = (100*weight/power(10,-6.287+3.330*log(10,length)))
      where species = 'SNSG'
      and nvl(length,0) > 119 and nvl(weight,0) <> 0
      and upload_session_id = p_upload_session_id;
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
        and upload_session_id = p_upload_session_id;        
            
 ----- Count Fish records uploaded to DS_FISH
 -- count how many new mr records were added to DS_MORIVER
  select count(*)
   into p_fishCntFinal
  from ds_fish
  where upload_session_id = p_upload_session_id; 
  
 -------------------  SUPPLEMENTAL   ---------------------------------------
 
 -- assign corresponding mr_datashee mr_id to supp records mr_id
 update ds_supplemental_check
 set (mr_id,complete,checkby) = 
     (SELECT mr_id, complete, checkby
     FROM ds_moriver
     WHERE ds_supplemental_check.MR_FID = ds_moriver.MR_FID
      and ds_moriver.UPLOAD_SESSION_ID = p_upload_session_id ),
      
     f_id = 
     (SELECT f_id 
     FROM ds_fish
     WHERE ds_supplemental_check.F_FID = ds_fish.F_FID
      and ds_fish.UPLOAD_SESSION_ID = p_upload_session_id )
      
 WHERE upload_session_id = p_upload_session_id;
 
  p_debug ('Around line 2663. p_suppSessionID = ' || p_upload_session_id, NULL);

  -- IF adding supplemental that was not in first upload - everything else the same
  FOR Y in (
    select mr_id, complete, checkby
    from ds_supplemental_check
    where upload_session_id = p_upload_session_id)
  
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
      
     WHERE upload_session_id = p_upload_session_id;
 		 
 	END IF; 
      
  END LOOP;
  p_debug ('Around line 2692. p_suppSessionID = ' || p_upload_session_id, NULL);

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
         and upload_session_id = p_upload_session_id
         and EXISTS (select f_fid from ds_fish
                         where ds_fish.f_fid = ds_supplemental_check.f_fid);   -- added to only add supplemental where matching MR record ;
         
  p_debug ('processFinalUpload Getting Count of Supplemetnal Records', NULL);
            
  ----- Count Supplemental records uploaded to DS_SUPPLEMENAL
  select count(*)
   into p_suppCntFinal
  from ds_supplemental
  where upload_session_id = p_upload_session_id;
  
  IF p_suppCntFinal > 0 THEN
      
    FOR Z in (
    select s_id, head, inter, m_ib, l_ob, l_ib, r_ib, r_ob, anal, dorsal
    from ds_supplemental
    where upload_session_id = p_upload_session_id
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
          
  END IF;
  
  -- insert into ds_procedure
  INSERT INTO ds_procedure (ID, F_FID, PURPOSE_CODE, PROCEDURE_DATE, PROCEDURE_START_TIME, PROCEDURE_END_TIME, PROCEDURE_BY, 
                                                    ANTIBIOTIC_INJECTION_IND, PHOTO_DORSAL_IND, PHOTO_VENTRAL_IND, PHOTO_LEFT_IND,
                                                    OLD_RADIO_TAG_NUM, OLD_FREQUENCY_ID, DST_SERIAL_NUM, DST_START_DATE, DST_START_TIME, DST_REIMPLANT_IND, NEW_RADIO_TAG_NUM,
                                                    NEW_FREQUENCY_ID, SEX_CODE, BLOOD_SAMPLE_IND, EGG_SAMPLE_IND, COMMENTS, FISH_HEALTH_COMMENTS,
                                                    EVAL_LOCATION_CODE, SPAWN_CODE, VISUAL_REPRO_STATUS_CODE, ULTRASOUND_REPRO_STATUS_CODE,
                                                    EXPECTED_SPAWN_YEAR, ULTRASOUND_GONAD_LENGTH, GONAD_CONDITION,
                                                    LAST_UPDATED, UPLOAD_SESSION_ID, UPLOADED_BY, UPLOAD_FILENAME )
    SELECT ID, F_FID, PURPOSE_CODE, PROCEDURE_DATE, PROCEDURE_START_TIME, PROCEDURE_END_TIME, PROCEDURE_BY, 
                                                    ANTIBIOTIC_INJECTION_IND, PHOTO_DORSAL_IND, PHOTO_VENTRAL_IND, PHOTO_LEFT_IND,
                                                    OLD_RADIO_TAG_NUM, OLD_FREQUENCY_ID, DST_SERIAL_NUM, DST_START_DATE, DST_START_TIME, DST_REIMPLANT_IND, NEW_RADIO_TAG_NUM,
                                                    NEW_FREQUENCY_ID, SEX_CODE, BLOOD_SAMPLE_IND, EGG_SAMPLE_IND, COMMENTS, FISH_HEALTH_COMMENTS,
                                                    EVAL_LOCATION_CODE, SPAWN_CODE, VISUAL_REPRO_STATUS_CODE, ULTRASOUND_REPRO_STATUS_CODE,
                                                    EXPECTED_SPAWN_YEAR, ULTRASOUND_GONAD_LENGTH, GONAD_CONDITION,
                                                    LAST_UPDATED, UPLOAD_SESSION_ID, UPLOADED_BY, UPLOAD_FILENAME 
    FROM UPLOAD_PROCEDURE up
    WHERE upload_session_id = p_upload_session_id
      AND up.f_fid NOT IN (SELECT f_fid FROM ds_procedure);   
  UPDATE ds_procedure dp SET f_id = (SELECT f_id FROM ds_fish df WHERE dp.f_fid = df.f_fid);

  SELECT COUNT(*) INTO p_procedureCntFinal FROM ds_procedure WHERE upload_session_id = p_upload_session_id;   

  COMMIT;
          
  EXCEPTION
  
    WHEN noSiteID then
        p_noSiteID_msg := ' '||p_noSite_cnt||' Missouri River datasheets are not associated with an existing or created site.  If new sites were created in the field, please upload PSPA_SITES_DATASHEET...CSV along with the MR, Fish and Supplemental datasheets.';
    WHEN OTHERS THEN    
      raise_application_error(-20001,'An error was encountered - '||SQLCODE||' -ERROR- '||SQLERRM);
      
-- l_session := 0;
-- l_siteMatch := 0;
-- l_siteMatch_tot := 0;
  
  END uploadFinal;

END data_upload_MAR2020;