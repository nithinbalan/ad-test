-- +goose Up
-- +goose StatementBegin

ALTER TABLE risk_register
ADD sub_section_id INT NULL;

ALTER TABLE risk_register
ADD CONSTRAINT FK_risk_register_sub_section_id
FOREIGN KEY (sub_section_id) REFERENCES organization_units(id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE risk_register
DROP CONSTRAINT FK_risk_register_sub_section_id;

ALTER TABLE risk_register
DROP COLUMN sub_section_id;

-- +goose StatementEnd