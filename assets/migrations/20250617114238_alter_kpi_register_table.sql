-- +goose Up
-- +goose StatementBegin
ALTER TABLE kpi_register
ADD entry_type_id INT NULL;

ALTER TABLE kpi_register
ADD data_source_id INT NULL;

ALTER TABLE kpi_register
ADD last_reviewed_date DATE NULL;

ALTER TABLE kpi_register
ADD next_review_date DATE NULL;

ALTER TABLE kpi_register
ADD last_review_note NVARCHAR(MAX) NULL;

ALTER TABLE kpi_register
ADD CONSTRAINT FK_kpi_register_entry_types
FOREIGN KEY (entry_type_id) REFERENCES entry_types(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE kpi_register
ADD CONSTRAINT FK_kpi_register_data_sources
FOREIGN KEY (data_source_id) REFERENCES data_sources(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE kpi_register
DROP CONSTRAINT FK_kpi_register_entry_types;

ALTER TABLE kpi_register
DROP COLUMN entry_type_id;

ALTER TABLE kpi_register
DROP CONSTRAINT FK_kpi_register_data_sources;

ALTER TABLE kpi_register
DROP COLUMN data_source_id;

ALTER TABLE kpi_register
DROP COLUMN last_reviewed_date;

ALTER TABLE kpi_register
DROP COLUMN next_review_date;

ALTER TABLE kpi_register
DROP COLUMN last_review_note;
-- +goose StatementEnd
