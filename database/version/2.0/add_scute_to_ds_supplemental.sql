--Add new fields to store scute data.
--Data will be migrated from the existing scute fields in the DS_SUPPLEMENTAL table.
alter table DS_SUPPLEMENTAL add (
  lscute number(1),
  rscute number(1),
  dscute number(1)
);