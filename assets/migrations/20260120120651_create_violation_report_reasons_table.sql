-- +goose Up
-- +goose StatementBegin
CREATE TABLE violation_report_reasons (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_violation_report_reasons_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.violation_report_reasons;
-- +goose StatementEnd
