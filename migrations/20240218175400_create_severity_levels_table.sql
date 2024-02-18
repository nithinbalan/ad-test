-- +goose Up
-- +goose StatementBegin
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = 'dbo' AND TABLE_NAME = 'severity_levels')
BEGIN
    CREATE TABLE dbo.severity_levels (
        id INT PRIMARY KEY IDENTITY(1, 1),
        title_en NVARCHAR(45) NOT NULL,
        title_ar NVARCHAR(45) NOT NULL,
        color_code NVARCHAR(45) NOT NULL,
        type NVARCHAR(10) NOT NULL,
        CONSTRAINT CHK_severity_levels_type CHECK (type IN ('high', 'medium', 'low'))
    )
END

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.severity_levels;
-- +goose StatementEnd
