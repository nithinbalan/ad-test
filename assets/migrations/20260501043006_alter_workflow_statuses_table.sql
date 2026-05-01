-- +goose Up
-- +goose StatementBegin

ALTER TABLE workflow_statuses
DROP CONSTRAINT CHK_workflow_statuses_type;

ALTER TABLE workflow_statuses
ALTER COLUMN type VARCHAR(25);

ALTER TABLE workflow_statuses
ADD CONSTRAINT CHK_workflow_statuses_type
CHECK (
    type IN (
        'submit',
        'revert',
        'reject',
        'approve',
        'publish',
        'pending',
        'comment',
        'changes_requested'
    )
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE workflow_statuses
DROP CONSTRAINT CHK_workflow_statuses_type;

ALTER TABLE workflow_statuses
ALTER COLUMN type VARCHAR(10);

ALTER TABLE workflow_statuses
ADD CONSTRAINT CHK_workflow_statuses_type
CHECK (
    type IN (
        'submit',
        'revert',
        'reject',
        'approve',
        'publish',
        'pending',
        'comment'
    )
);

-- +goose StatementEnd
