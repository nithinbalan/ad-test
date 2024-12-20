-- +goose Up
-- +goose StatementBegin
-- Add the recovery objective columns to the documents table
ALTER TABLE documents 
ADD rto INT NULL;

ALTER TABLE documents 
ADD rto_unit_id INT NULL;

ALTER TABLE documents 
ADD mtpd INT NULL;

ALTER TABLE documents 
ADD mtpd_unit_id INT NULL;

ALTER TABLE documents 
ADD rpo INT NULL;

ALTER TABLE documents 
ADD rpo_unit_id INT NULL;

-- Add the foreign key constraints for rto_unit_id, mtpd_unit_id, and rpo_unit_id
ALTER TABLE documents
ADD CONSTRAINT FK_documents_rto_time_units FOREIGN KEY (rto_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE documents
ADD CONSTRAINT FK_documents_mtpd_time_units FOREIGN KEY (mtpd_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE documents
ADD CONSTRAINT FK_documents_rpo_time_units FOREIGN KEY (rpo_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove foreign key constraints before dropping columns
ALTER TABLE documents
DROP CONSTRAINT FK_documents_rto_time_units;

ALTER TABLE documents
DROP COLUMN rto_unit_id;

ALTER TABLE documents
DROP CONSTRAINT FK_documents_mtpd_time_units;

ALTER TABLE documents
DROP COLUMN mtpd_unit_id;

ALTER TABLE documents
DROP CONSTRAINT FK_documents_rpo_time_units;

ALTER TABLE documents
DROP COLUMN rpo_unit_id;

-- Remove the additional columns
ALTER TABLE documents 
DROP COLUMN rto;

ALTER TABLE documents 
DROP COLUMN mtpd;

ALTER TABLE documents 
DROP COLUMN rpo;
-- +goose StatementEnd
