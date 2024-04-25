-- +goose Up
-- +goose StatementBegin
CREATE TABLE severity_levels (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(45) NOT NULL,
    title_ar NVARCHAR(45) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(10) NOT NULL,
        CONSTRAINT CHK_severity_levels_type CHECK (type IN ('high', 'medium', 'low')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_severity_levels_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.severity_levels;
-- +goose StatementEnd
