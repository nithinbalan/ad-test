-- +goose Up
-- +goose StatementBegin

ALTER TABLE incident_statuses
DROP CONSTRAINT CHK_incident_statuses_type;

ALTER TABLE incident_statuses
ADD CONSTRAINT CHK_incident_statuses_type
CHECK (
    type IN (
        'open',
        'investigating',
        'in_progress',
        'resolved',
        'closed',
        'pending_review',
        'escalated',
        're_opened',
        'deferred',
        'mitigated',
        'new',
        'triage',
        'containment',
        'eradication',
        'recovery',
        'closure_preparation'
    )
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE incident_statuses
DROP CONSTRAINT CHK_incident_statuses_type;

ALTER TABLE incident_statuses
ADD CONSTRAINT CHK_incident_statuses_type
CHECK (
    type IN (
        'open',
        'investigating',
        'in_progress',
        'resolved',
        'closed',
        'pending_review',
        'escalated',
        're_opened',
        'deferred',
        'mitigated'
    )
);

-- +goose StatementEnd
