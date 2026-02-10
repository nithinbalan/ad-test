-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_register
ADD note NVARCHAR(MAX) NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_register
DROP COLUMN note;
-- +goose StatementEnd