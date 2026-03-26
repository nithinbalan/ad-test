-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_register_issuing_authority_entities
ALTER COLUMN hijri_date_of_receiving NVARCHAR(500) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_register_issuing_authority_entities
ALTER COLUMN hijri_date_of_receiving DATETIME NOT NULL;
-- +goose StatementEnd
