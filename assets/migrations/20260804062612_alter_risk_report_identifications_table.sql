-- +goose Up
-- +goose StatementBegin
ALTER TABLE risk_report_identifications
ADD registration_date DATE NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE risk_report_identifications
DROP COLUMN registration_date;
-- +goose StatementEnd