-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_response_strategies (
    id INT PRIMARY KEY IDENTITY(1, 1),
    type VARCHAR(25) NOT NULL,
        CONSTRAINT CHK_risk_response_strategies_type CHECK (type IN ('accept', 'avoid', 'transfer', 'mitigate', 'tolerate')),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    description_en NVARCHAR(MAX),
    description_ar NVARCHAR(MAX),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_risk_response_strategies_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.risk_response_strategies;
-- +goose StatementEnd
