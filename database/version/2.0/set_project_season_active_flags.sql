ALTER TABLE project_lk
ADD ACTIVE_FLAG_TF VARCHAR2(1 BYTE) NULL;

UPDATE project_lk
SET ACTIVE_FLAG_TF = 'T';

-- --------------------------------------

ALTER TABLE season_lk
ADD ACTIVE_FLAG_TF VARCHAR2(1 BYTE) NULL;

UPDATE season_lk
SET ACTIVE_FLAG_TF = 'T';