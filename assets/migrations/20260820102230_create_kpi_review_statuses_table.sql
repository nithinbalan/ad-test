-- +goose Up
-- +goose StatementBegin
CREATE TABLE kpi_review_statuses (
    id INT IDENTITY(1, 1) PRIMARY KEY,
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    color_code VARCHAR(20) NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(10) NOT NULL CONSTRAINT DF_kpi_review_statuses_status DEFAULT 'active',
    CONSTRAINT CHK_kpi_review_statuses_type CHECK (
        [type] IN ('updated', 'awaiting-update', 'overdue', 'submitted', 'approved', 'reverted')
    ),
    CONSTRAINT CHK_kpi_review_statuses_status CHECK ([status] IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.kpi_review_statuses;
-- +goose StatementEnd
