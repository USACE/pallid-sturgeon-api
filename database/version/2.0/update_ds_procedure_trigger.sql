--Undo...
/*
create or replace TRIGGER bi_DS_PROCEDURE
  before insert on DS_PROCEDURE            
  for each row 
begin  
  if :new.ID is null then
    select PROCEDURE_ID_SEQ.nextval into :new.ID from dual;
  end if;
end;
*/

-- Updates trigger to use the more accurate "PROCEDURE_SEQ" sequence instead of
-- "PROCEDURE_ID_SEQ" (whose LAST_NUMBER is much lower than the max(ds_procedure.id)).
-- PROCEDURE_SEQ appears to be the correct sequence to use for the DS_PROCEDURE table.

--Before running, please run the following to compare the max(id) in the table to the last number in the sequence:
--select max(id) from ds_procedure;  --in dev: 2355
--select sequence_name, last_number from user_sequences where sequence_name like 'PROCEDURE%';
-- in dev:
--PROCEDURE_ID_SEQ	320  -- way too low. probably was not used.
--PROCEDURE_SEQ	2380  -- looks correct.

create or replace TRIGGER bi_DS_PROCEDURE
  before insert on DS_PROCEDURE            
  for each row 
begin  
  if :new.ID is null then
    select PROCEDURE_SEQ.nextval into :new.ID from dual;
  end if;
end;

--Trigger BI_DS_PROCEDURE compiled