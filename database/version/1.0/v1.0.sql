
CREATE SEQUENCE procedure_seq NOCACHE;

@procedures_tables.sql;

@data_upload_mar2020.pks;
@data_upload_mar2020.pkb;



### Already run in Prod
ALTER SEQUENCE upload_session_seq INCREMENT BY 5000;
SELECT upload_session_seq.nextval from dual;
ALTER SEQUENCE upload_session_seq INCREMENT BY 1;