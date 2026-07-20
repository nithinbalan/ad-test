-- +goose Up
-- +goose StatementBegin

ALTER TABLE workflow_statuses
DROP CONSTRAINT CHK_workflow_statuses_type;

ALTER TABLE workflow_statuses
ADD CONSTRAINT CHK_workflow_statuses_type
CHECK (
    [type] IN (
        'submit',
        'revert',
        'reject',
        'approve',
        'publish',
        'pending',
        'comment',
        'reopen',
        'resubmit'
    )
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE workflow_statuses
DROP CONSTRAINT CHK_workflow_statuses_type;

DELETE FROM workflow_statuses
WHERE type = 'resubmit';

ALTER TABLE workflow_statuses
ADD CONSTRAINT CHK_workflow_statuses_type
CHECK (
    [type] IN (
        'submit',
        'revert',
        'reject',
        'approve',
        'publish',
        'pending',
        'comment',
        'reopen'
    )
);

-- +goose StatementEnd
