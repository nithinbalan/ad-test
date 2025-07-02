-- +goose Up
-- +goose StatementBegin
ALTER TABLE workflows
ADD module_identifier NVARCHAR(255) NULL;

ALTER TABLE workflows
ADD title_en NVARCHAR(255) NULL;

ALTER TABLE workflows
ADD title_ar NVARCHAR(255) NULL;

ALTER TABLE workflows
DROP CONSTRAINT CHK_workflows_type;

ALTER TABLE workflows
ADD CONSTRAINT CHK_workflows_type CHECK (type IN ('document', 'draft-report', 'final-report'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workflows
DROP COLUMN module_identifier;

ALTER TABLE workflows
DROP COLUMN title_en;

ALTER TABLE workflows
DROP COLUMN title_ar;

ALTER TABLE workflows
DROP CONSTRAINT CHK_workflows_type;

ALTER TABLE workflows
ADD CONSTRAINT CHK_workflows_type CHECK (type IN ('document'));
-- +goose StatementEnd
