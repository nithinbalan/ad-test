-- +goose Up
-- +goose StatementBegin
CREATE TABLE kpi_achievement_logics (
    id INT IDENTITY(1, 1) PRIMARY KEY,
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    status VARCHAR(10) NOT NULL CONSTRAINT DF_kpi_achievement_logics_status DEFAULT 'active',
    type VARCHAR(25) NOT NULL,
    CONSTRAINT CHK_kpi_achievement_logics_status CHECK ([status] IN ('active', 'inactive')),
    CONSTRAINT CHK_kpi_achievement_logics_type CHECK (
        [type] IN (
            'lower', 'higher', 'steady_or_increasing',
            'increasing', 'decreasing', 'steady_or_decreasing'
        )
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.kpi_achievement_logics;
-- +goose StatementEnd
