-- +goose Up
-- +goose StatementBegin

CREATE TABLE violation_report_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(100) NOT NULL,
        CONSTRAINT CHK_violation_report_statuses_type CHECK (
            type IN (
                'new',
                'under_investigation',
                'needs_clarification',
                'in_progress',
                'under_review',
                'on_hold',
                'closed_rejected',
                'closed_resolved',
                'closed_irrelevant',
                'closed_no_violation_found',
                'closed_referred_to_entity'
            )
        ),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_violation_report_statuses_status CHECK (status IN ('active', 'inactive'))
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE dbo.violation_report_statuses;

-- +goose StatementEnd
