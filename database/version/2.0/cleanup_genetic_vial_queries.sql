-- Insert/Update Statements

insert into ds_supplemental (s_id, fieldoffice, project, segment, genetics_vial_number, f_id, mr_id, uploaded_by, f_fid, last_updated)
select supplemental_id_seq.nextval, fieldoffice, project, segment, genetics_vial_number, f_id, mr_id, uploaded_by, f_fid, sysdate
from ds_fish
where genetics_vial_number is not null and f_id not in (2193550, 2350281, 2350285)  and (f_id, mr_id) not in (select f_id, mr_id from ds_supplemental);

update ds_supplemental
set genetics_vial_number = '22-00576'
where f_id = 2198872;
Select Statements
select f_id, mr_id, f_fid, mr_fid, species, genetics_vial_number from ds_fish
where genetics_vial_number is not null;

select f_id, mr_id, f_fid, mr_fid, species, genetics_vial_number from ds_fish
where genetics_vial_number is not null and (f_id, mr_id) not in (select f_id, mr_id from ds_supplemental);

select count(*) from ds_fish
where genetics_vial_number is not null;

select ds_fish.f_id, ds_fish.mr_id, ds_fish.f_fid, ds_fish.species, ds_fish.genetics_vial_number, 
ds_supplemental.f_id, ds_supplemental.mr_id, ds_supplemental.genetics_vial_number 
from ds_fish, ds_supplemental
where ds_fish.genetics_vial_number is not null and ds_fish.f_id = ds_supplemental.f_id and ds_fish.mr_id = ds_supplemental.mr_id
and ds_supplemental.genetics_vial_number is null;

select count(*)
from ds_fish, ds_supplemental
where ds_fish.genetics_vial_number is not null and ds_fish.f_id = ds_supplemental.f_id and ds_fish.mr_id = ds_supplemental.mr_id;

select f_id, mr_id, genetics_vial_number from ds_supplemental where genetics_vial_number is not null;

select count(*) from ds_supplemental where genetics_vial_number is not null;

select f_id, mr_id, f_fid, genetics_vial_number, tagnumber
from ds_supplemental
where tagnumber is not null;

select f_id, mr_id, last_updated, uploaded_by, upload_filename, genetics_vial_number, f_fid
from ds_supplemental
where (f_id, mr_id) not in (select f_id, mr_id from ds_fish);

select * from ds_procedure 
where (f_id) not in (select f_id from ds_fish);

select * from ds_fish where f_fid = '20240318-141843836-004-003';

select * from ds_fish where f_id in ('2437995', '2437998', '2280613');

select * from ds_procedure where f_id = '2280564';

select * from ds_procedure where f_id not in (select f_id from ds_supplemental);
