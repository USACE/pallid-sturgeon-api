-- NEW PROCEDURE FIELDS
ALTER TABLE ds_procedure
ADD (
    old_rt_serial NUMBER,
    old_dst_serial NUMBER,
    new_rt_serial NUMBER,
    polarization_index NUMBER(3,2)
);