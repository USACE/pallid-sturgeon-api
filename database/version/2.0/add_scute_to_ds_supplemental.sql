--Add new fields to store scute data.
--Data will be migrated from the existing scute fields in the DS_SUPPLEMENTAL table.
alter table DS_SUPPLEMENTAL add (
  lscute number(1),
  rscute number(1),
  dscute number(1)
);


update ds_supplemental
set lscute = scutenum
where scuteloc = 'L';

update ds_supplemental
set lscute = scutenum2
where scuteloc2 = 'L' and lscute is null;

update ds_supplemental
set rscute = scutenum
where scuteloc = 'R';

update ds_supplemental
set rscute = scutenum2
where scuteloc2 = 'R' and rscute is null;

update ds_supplemental
set dscute = scutenum
where scuteloc = 'D';

update ds_supplemental
set dscute = scutenum2
where scuteloc2 = 'D' and dscute is null;