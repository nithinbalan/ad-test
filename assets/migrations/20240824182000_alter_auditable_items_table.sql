-- +goose Up
-- +goose StatementBegin
ALTER TABLE auditable_items
ADD type NVARCHAR(50) NOT NULL DEFAULT 'custom',
CONSTRAINT CHK_auditable_items_type CHECK (type IN ('process', 'control', 'incident', 'risk', 'document', 'compliance_register', 'compliance_requirement', 'user', 'custom'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE auditable_items
DROP CONSTRAINT CHK_auditable_items_type;
ALTER TABLE auditable_items
DROP COLUMN type;
-- +goose StatementEnd
