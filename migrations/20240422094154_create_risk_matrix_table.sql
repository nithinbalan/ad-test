-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_matrix (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title NVARCHAR(100) NOT NULL,
    calculation_type VARCHAR(20) NOT NULL DEFAULT 'multiplication',
    type VARCHAR(10) NOT NULL DEFAULT 'erm',
    CONSTRAINT CHK_risk_matrix_calculation_type CHECK (calculation_type IN ('multiplication', 'addition')),
    CONSTRAINT CHK_risk_matrix_type CHECK (type IN ('erm', 'isms'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_matrix;
-- +goose StatementEnd
