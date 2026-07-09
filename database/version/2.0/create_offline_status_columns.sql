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

-- 2. add status column to ds_moriver

-- add status column
alter table ds_moriver
add status number(1);

-- fill existing data for status column
update ds_moriver
set status = 2
where status is null;

-- alter DATA_ENTRY_MISSOURI_OBJ_TYPE object type
ALTER TYPE DATA_ENTRY_MISSOURI_OBJ_TYPE ADD ATTRIBUTE (status NUMBER(1)) CASCADE

-- @TODO: Manually had status column in package pallid_data_entry_api.data_entry_missouri_fnc

commit;