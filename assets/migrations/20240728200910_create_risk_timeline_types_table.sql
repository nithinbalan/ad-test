-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_timeline_types
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL DEFAULT 'identified',
    CONSTRAINT CHK_risk_timeline_types_type CHECK (type IN (
        'identified', 
        'rca', 
        'risk_analysis', 
        'control_analysis', 
        'risk_response', 
        'residual_control_analysis', 
        'residual_risk_analysis', 
        'monitoring', 
        'closure'
    ))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.risk_timeline_types;
-- +goose StatementEnd
