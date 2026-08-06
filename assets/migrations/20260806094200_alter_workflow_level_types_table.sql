-- +goose Up
-- +goose StatementBegin

ALTER TABLE workflow_level_types
ALTER COLUMN title_en VARCHAR(500) NOT NULL;

ALTER TABLE workflow_level_types
ALTER COLUMN title_ar VARCHAR(500) NOT NULL;

ALTER TABLE workflow_level_types
ALTER COLUMN type VARCHAR(500) NOT NULL;

ALTER TABLE workflow_level_types
DROP CONSTRAINT CHK_workflow_level_types_type;

ALTER TABLE workflow_level_types
ADD CONSTRAINT CHK_workflow_level_types_type
CHECK (
    type IN (
        'individual',
        'multiple-users',
        'role',
        'hod',
        'custom-roles',
        'user-and-reporting-user'
    )
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE workflow_level_types
DROP CONSTRAINT CHK_workflow_level_types_type;

ALTER TABLE workflow_level_types
ADD CONSTRAINT CHK_workflow_level_types_type
CHECK (
    type IN (
        'individual',
        'multiple-users',
        'role',
        'hod',
        'custom-roles'
    )
);

ALTER TABLE workflow_level_types
ALTER COLUMN type VARCHAR(20) NOT NULL;

-- +goose StatementEnd