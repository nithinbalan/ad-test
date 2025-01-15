-- +goose Up
-- +goose StatementBegin
ALTER TABLE incident_register
ALTER COLUMN incident_type_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE incident_register
ALTER COLUMN incident_type_id INT NOT NULL;
-- +goose StatementEnd
