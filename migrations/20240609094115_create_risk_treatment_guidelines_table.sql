-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_treatment_guidelines (
    id INT PRIMARY KEY IDENTITY(1, 1),
    risk_matrix_risk_rating_id INT NOT NULL,
    description_en NVARCHAR(100) NOT NULL,
    description_ar NVARCHAR(100) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',

    CONSTRAINT CHK_risk_treatment_guidelines_status CHECK (status IN ('active', 'inactive')),
    CONSTRAINT FK_risk_treatment_guidelines_risk_matrix_risk_ratings FOREIGN KEY (risk_matrix_risk_rating_id) REFERENCES risk_matrix_risk_ratings (id) ON DELETE NO ACTION ON UPDATE NO ACTION

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.risk_treatment_guidelines;
-- +goose StatementEnd
