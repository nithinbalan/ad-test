-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_treatment_guidelines (
    id INT PRIMARY KEY IDENTITY(1, 1),
    risk_level_id INT NOT NULL,
    timeframe_en NVARCHAR(100) NOT NULL,
    timeframe_ar NVARCHAR(100) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_risk_treatment_guidelines_status CHECK (status IN ('active', 'inactive')),
    CONSTRAINT FK_risk_treatment_guidelines_risk_levels FOREIGN KEY (risk_level_id) REFERENCES risk_levels (id) ON DELETE NO ACTION ON UPDATE NO ACTION

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.risk_treatment_guidelines;
-- +goose StatementEnd
