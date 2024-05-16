-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_evaluations (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_register_id INT NOT NULL,
    inherent_risk_matrix_likelihood_id INT NOT NULL,
    inherent_risk_matrix_impact_id INT NOT NULL,
    inherent_risk_matrix_risk_rating_id INT NOT NULL,
    residual_risk_matrix_likelihood_id INT NULL,
    residual_risk_matrix_impact_id INT NULL,
    residual_risk_matrix_risk_rating_id INT NULL,
    created_by INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    updated_by INT NULL,
    CONSTRAINT FK_risk_evaluation_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register(id),
    CONSTRAINT FK_risk_evaluation_risk_matrix_likelihoods_1 FOREIGN KEY (inherent_risk_matrix_likelihood_id) REFERENCES risk_matrix_likelihoods(id),
    CONSTRAINT FK_risk_evaluation_risk_matrix_impacts_1 FOREIGN KEY (inherent_risk_matrix_impact_id) REFERENCES risk_matrix_impacts(id),
    CONSTRAINT FK_risk_evaluation_risk_matrix_risk_ratings_1 FOREIGN KEY (inherent_risk_matrix_risk_rating_id) REFERENCES risk_matrix_risk_ratings(id),
    CONSTRAINT FK_risk_evaluation_risk_matrix_likelihoods_2 FOREIGN KEY (residual_risk_matrix_likelihood_id) REFERENCES risk_matrix_likelihoods(id),
    CONSTRAINT FK_risk_evaluation_risk_matrix_impacts_2 FOREIGN KEY (residual_risk_matrix_impact_id) REFERENCES risk_matrix_impacts(id),
    CONSTRAINT FK_risk_evaluation_risk_matrix_risk_ratings_2 FOREIGN KEY (residual_risk_matrix_risk_rating_id) REFERENCES risk_matrix_risk_ratings(id),
    CONSTRAINT FK_risk_evaluation_users_1 FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT FK_risk_evaluation_users_2 FOREIGN KEY (updated_by) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS risk_evaluations;
-- +goose StatementEnd
