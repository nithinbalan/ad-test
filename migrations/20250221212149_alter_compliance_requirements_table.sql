-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_requirements
    ALTER COLUMN compliance_requirement_category_id INT NULL;

ALTER TABLE compliance_requirements
    ALTER COLUMN severity_level_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_requirements
    ALTER COLUMN compliance_requirement_category_id INT NOT NULL;

ALTER TABLE compliance_requirements
    ALTER COLUMN severity_level_id INT NOT NULL;
-- +goose StatementEnd