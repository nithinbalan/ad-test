-- +goose Up
-- +goose StatementBegin
CREATE TABLE process_criticality_levels (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL,
        CONSTRAINT CHK_process_criticality_levels_type CHECK (type IN ('high', 'medium','low')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_process_criticality_levels_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.process_criticality_levels;
-- +goose StatementEnd

