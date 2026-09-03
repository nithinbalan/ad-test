-- +goose Up
-- +goose StatementBegin

CREATE TABLE planned_advisory_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    description_en NVARCHAR(MAX) NOT NULL,
    description_ar NVARCHAR(MAX) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(50) NOT NULL,
        CONSTRAINT CHK_planned_advisory_statuses_type CHECK (
            type IN (
                'open',
                'acknowledged',
                'awaits_feedback',
                'feedback_under_review',
                'awaits_feedback_escalated',
                'resolved',
                'rejected',
                'reopened'
            )
        ),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_planned_advisory_statuses_status CHECK (status IN ('active', 'inactive'))
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE planned_advisory_statuses;

-- +goose StatementEnd

