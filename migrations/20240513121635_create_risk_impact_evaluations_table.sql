-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_impact_evaluations (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_register_id INT NOT NULL,
    risk_matrix_impact_guideline_id INT NOT NULL,
    risk_impact_category_id INT NOT NULL,
    risk_matrix_impact_id INT NOT NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_risk_impact_evaluation_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register(id),
    CONSTRAINT FK_risk_impact_evaluation_risk_impact_category FOREIGN KEY (risk_impact_category_id) REFERENCES risk_impact_categories(id),
    CONSTRAINT FK_risk_impact_evaluation_risk_matrix_impacts FOREIGN KEY (risk_matrix_impact_id) REFERENCES risk_matrix_impacts(id),
    CONSTRAINT FK_risk_impact_evaluation_risk_matrix_impact_guidelines FOREIGN KEY (risk_matrix_impact_guideline_id) REFERENCES risk_matrix_impact_guidelines(id),
    CONSTRAINT FK_risk_impact_evaluation_users_1 FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT FK_risk_impact_evaluation_users_2 FOREIGN KEY (updated_by) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS risk_impact_evaluations;
-- +goose StatementEnd
