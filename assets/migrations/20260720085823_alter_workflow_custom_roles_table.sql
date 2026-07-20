-- +goose Up
-- +goose StatementBegin

ALTER TABLE workflow_custom_roles
DROP CONSTRAINT CHK_workflow_custom_roles_type;

ALTER TABLE workflow_custom_roles
ADD CONSTRAINT CHK_workflow_custom_roles_type
CHECK (
    [type] IN (
        'department_head',
        'risk_responsible_user'
    )
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE workflow_custom_roles
DROP CONSTRAINT CHK_workflow_custom_roles_type;

ALTER TABLE workflow_custom_roles
ADD CONSTRAINT CHK_workflow_custom_roles_type
CHECK (
    [type] = 'department_head'
);

-- +goose StatementEnd
