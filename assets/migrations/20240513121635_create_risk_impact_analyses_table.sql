-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_impact_analyses (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_evaluation_id INT NOT NULL,
    risk_impact_category_id INT NOT NULL,
    risk_matrix_impact_id INT NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_risk_impact_analyses_risk_evaluations FOREIGN KEY (risk_evaluation_id) REFERENCES risk_evaluations(id),
    CONSTRAINT FK_risk_impact_analyses_risk_impact_categories FOREIGN KEY (risk_impact_category_id) REFERENCES risk_impact_categories(id),
    CONSTRAINT FK_risk_impact_analyses_risk_matrix_impacts FOREIGN KEY (risk_matrix_impact_id) REFERENCES risk_matrix_impacts(id),
    CONSTRAINT FK_risk_impact_analyses_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id),
    CONSTRAINT FK_risk_impact_analyses_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS risk_impact_analyses;
-- +goose StatementEnd
