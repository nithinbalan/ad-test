-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_evaluations (
    id INT PRIMARY KEY IDENTITY(1,1),
    type VARCHAR(10) NOT NULL DEFAULT 'inherent',
        CONSTRAINT CHK_risk_evaluations_type CHECK (type IN ('inherent', 'residual')),
    risk_register_id INT NOT NULL,
    risk_matrix_likelihood_id INT NOT NULL,
    risk_matrix_impact_id INT NOT NULL,
    risk_matrix_risk_rating_id INT NOT NULL,
    risk_matrix_risk_rating_value_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_risk_evaluations_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register(id),
    CONSTRAINT FK_risk_evaluations_risk_matrix_likelihoods_1 FOREIGN KEY (risk_matrix_likelihood_id) REFERENCES risk_matrix_likelihoods(id),
    CONSTRAINT FK_risk_evaluations_risk_matrix_impacts_1 FOREIGN KEY (risk_matrix_impact_id) REFERENCES risk_matrix_impacts(id),
    CONSTRAINT FK_risk_evaluations_risk_matrix_risk_ratings_1 FOREIGN KEY (risk_matrix_risk_rating_id) REFERENCES risk_matrix_risk_ratings(id),
    CONSTRAINT FK_risk_evaluations_risk_matrix_risk_rating_values FOREIGN KEY (risk_matrix_risk_rating_value_id) REFERENCES risk_matrix_risk_rating_values(id),
    CONSTRAINT FK_risk_evaluations_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id),
    CONSTRAINT FK_risk_evaluations_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS risk_evaluations;
-- +goose StatementEnd
