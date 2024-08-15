-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_timeline_types (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL DEFAULT 'reporting',
        CONSTRAINT CHK_incident_timeline_types_type CHECK (type IN ('reporting', 'investigating', 'rca', 'communication', 'action-plans', 'closure', 'pir'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.incident_timeline_types;
-- +goose StatementEnd
