ALTER TABLE miners
ADD COLUMN batch_label varchar(50) null;

ALTER TABLE miners
ADD COLUMN gifted BOOLEAN null;