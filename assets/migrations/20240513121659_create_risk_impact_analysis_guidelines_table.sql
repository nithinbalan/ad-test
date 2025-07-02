-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_impact_analysis_guidelines (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_impact_analysis_id INT NOT NULL,
    risk_matrix_impact_guideline_id INT NOT NULL,
    CONSTRAINT FK_risk_impact_analysis_guidelines_risk_impact_analyses FOREIGN KEY (risk_impact_analysis_id) REFERENCES risk_impact_analyses(id),
    CONSTRAINT FK_risk_impact_analysis_guidelines_risk_matrix_impact_guidelines FOREIGN KEY (risk_matrix_impact_guideline_id) REFERENCES risk_matrix_impact_guidelines(id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS risk_impact_analysis_guidelines;
-- +goose StatementEnd
