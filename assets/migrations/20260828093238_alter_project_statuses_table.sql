-- +goose Up
-- +goose StatementBegin
ALTER TABLE project_statuses
DROP CONSTRAINT CHK_project_statuses_type;

ALTER TABLE project_statuses
ADD CONSTRAINT CHK_project_statuses_type CHECK (
    type IN (
        'not_started',
        'planning',
        'in_progress',
        'on_hold',
        'cancelled',
        'completed',
        'on_track',
        'at_risk',
        'delayed'
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE project_statuses
DROP CONSTRAINT CHK_project_statuses_type;

ALTER TABLE project_statuses
ADD CONSTRAINT CHK_project_statuses_type CHECK (
    type IN (
        'not_started',
        'planning',
        'in_progress',
        'on_hold',
        'cancelled',
        'completed'
    )
);
-- +goose StatementEnd
