-- +goose Up
-- +goose StatementBegin

ALTER TABLE violation_reports
DROP CONSTRAINT FK_violation_reports_organization_units;

ALTER TABLE violation_reports
DROP COLUMN reported_entity_id;

ALTER TABLE violation_reports
ADD reported_entity NVARCHAR(MAX) NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE violation_reports
DROP COLUMN reported_entity;

ALTER TABLE violation_reports
ADD reported_entity_id INT NULL;

ALTER TABLE violation_reports
ADD CONSTRAINT FK_violation_reports_organization_units
FOREIGN KEY (reported_entity_id)
REFERENCES organization_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

-- +goose StatementEnd