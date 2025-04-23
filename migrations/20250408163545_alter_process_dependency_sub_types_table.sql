-- +goose Up
-- +goose StatementBegin
ALTER TABLE process_dependency_sub_types
ADD type VARCHAR(20) NOT NULL DEFAULT 'other',
    CONSTRAINT CHK_process_dependency_sub_types_type 
        CHECK (type IN ('process', 'department', 'asset', 'vendor', 'other'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE process_dependency_sub_types
DROP CONSTRAINT CHK_process_dependency_sub_types_type;

ALTER TABLE process_dependency_sub_types
DROP COLUMN type;
-- +goose StatementEnd
