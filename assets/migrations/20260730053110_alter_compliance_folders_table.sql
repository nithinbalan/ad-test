-- +goose Up
-- +goose StatementBegin

ALTER TABLE compliance_folders
ADD
    is_cyber BIT NOT NULL
        CONSTRAINT DF_compliance_folders_is_cyber DEFAULT (0);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE compliance_folders
DROP CONSTRAINT DF_compliance_folders_is_cyber;

ALTER TABLE compliance_folders
DROP COLUMN is_cyber;

-- +goose StatementEnd