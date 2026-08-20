-- +goose Up
-- +goose StatementBegin
ALTER TABLE kpi_statuses DROP CONSTRAINT CHK_kpi_statuses_type;
ALTER TABLE kpi_statuses ALTER COLUMN type VARCHAR(20) NOT NULL;

ALTER TABLE kpi_statuses
ADD label VARCHAR(255) NULL,
    is_not_achieved BIT NOT NULL
        CONSTRAINT DF_kpi_statuses_is_not_achieved DEFAULT 0,
    [order] INT NULL;

ALTER TABLE kpi_statuses
ADD CONSTRAINT CHK_kpi_statuses_type CHECK (
    [type] IN (
        'new', 'not-achieved', 'achieved', 'at-risk', 'on-track',
        'archived', 'inactive', 'active', 'low',
        'needs-improvement', 'near-target', 'exceeded'
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE kpi_statuses DROP CONSTRAINT CHK_kpi_statuses_type;
ALTER TABLE kpi_statuses DROP CONSTRAINT DF_kpi_statuses_is_not_achieved;
ALTER TABLE kpi_statuses DROP COLUMN label, is_not_achieved, [order];
ALTER TABLE kpi_statuses ALTER COLUMN type VARCHAR(10) NOT NULL;
ALTER TABLE kpi_statuses
ADD CONSTRAINT CHK_kpi_statuses_type
    CHECK ([type] IN ('active', 'inactive', 'archived'));
-- +goose StatementEnd
