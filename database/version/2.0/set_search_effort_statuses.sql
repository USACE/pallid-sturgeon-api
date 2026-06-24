Search Effort Form Validations Scripts

-- 1. add status column to ds_search

-- add status column
alter table ds_search
add status number(1);

-- fill existing data for status column
update ds_search
set status = 2
where status is null;

-- make stop time/latitude/longitude nullable
alter table ds_search
modify stop_time null;

alter table ds_search
modify stop_latitude null;

alter table ds_search
modify stop_longitude null; 

commit;