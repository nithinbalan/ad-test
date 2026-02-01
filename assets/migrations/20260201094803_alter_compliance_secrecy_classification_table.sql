-- +goose Up
-- +goose StatementBegin
-- Add the classification_type column
ALTER TABLE compliance_secrecy_classifications
ADD type VARCHAR(50) NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- Remove the classification_type column
ALTER TABLE compliance_secrecy_classifications
DROP COLUMN type;
-- +goose StatementEnd
