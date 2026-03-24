-- +goose Up
-- +goose StatementBegin

CREATE TABLE internal_committee_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(100) NOT NULL,
        CONSTRAINT CHK_internal_committee_statuses_type CHECK (
            type IN (
                'draft',
                'submitted',
                'in_review',
                'in_progress',
                'approved',
                'revert',
                'reformation',
                'disbanded'
            )
        ),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_internal_committee_statuses_status CHECK (status IN ('active', 'inactive'))
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE dbo.internal_committee_statuses;

-- +goose StatementEnd
