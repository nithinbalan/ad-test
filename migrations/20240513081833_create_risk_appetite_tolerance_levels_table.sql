-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_appetite_tolerance_levels (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(45) NOT NULL,
    title_ar NVARCHAR(45) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(20) NOT NULL,
        CONSTRAINT CHK_risk_appetite_tolerance_levels_type CHECK (type IN ('above', 'within', 'below', 'approaching', 'significantly_above', 'marginally_below')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_risk_appetite_tolerance_levels_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.risk_appetite_tolerance_levels;
-- +goose StatementEnd
