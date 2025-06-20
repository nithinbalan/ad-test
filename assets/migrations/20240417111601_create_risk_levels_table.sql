-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_levels (
    id INT PRIMARY KEY IDENTITY(1,1),
    type VARCHAR(10) NOT NULL,
        CONSTRAINT CHK_risk_levels_type CHECK (type IN ('high', 'medium', 'low')),
    title_en NVARCHAR(45) NOT NULL,
    title_ar NVARCHAR(45) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_risk_levels_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS risk_levels;
-- +goose StatementEnd
