insert into ds_supplemental (s_id, fieldoffice, project, segment, genetics_vial_number, f_id, mr_id, uploaded_by, f_fid, last_updated)
select supplemental_id_seq.nextval, fieldoffice, project, segment, genetics_vial_number, f_id, mr_id, uploaded_by, f_fid, sysdate
from ds_fish
where genetics_vial_number is not null and f_id not in (2193550, 2350281, 2350285)  and (f_id, mr_id) not in (select f_id, mr_id from ds_supplemental);