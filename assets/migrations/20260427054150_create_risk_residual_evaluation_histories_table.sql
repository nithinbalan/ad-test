-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_residual_evaluation_histories (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_evaluation_id INT NOT NULL,
    risk_register_id INT NOT NULL,
    risk_matrix_likelihood_id INT NOT NULL,
    risk_matrix_risk_rating_id INT NOT NULL,
    risk_matrix_risk_rating_value_id INT NOT NULL,
    average_control_effectiveness_score DECIMAL(5,2) DEFAULT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_risk_residual_evaluation_histories_risk_evaluation
        FOREIGN KEY (risk_evaluation_id)
        REFERENCES risk_evaluations(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_risk_residual_evaluation_histories_risk_register
        FOREIGN KEY (risk_register_id)
        REFERENCES risk_register(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_risk_residual_evaluation_histories_risk_matrix_likelihood
        FOREIGN KEY (risk_matrix_likelihood_id)
        REFERENCES risk_matrix_likelihoods(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_risk_residual_evaluation_histories_risk_matrix_risk_rating
        FOREIGN KEY (risk_matrix_risk_rating_id)
        REFERENCES risk_matrix_risk_ratings(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_risk_residual_evaluation_histories_risk_matrix_risk_rating_value
        FOREIGN KEY (risk_matrix_risk_rating_value_id)
        REFERENCES risk_matrix_risk_rating_values(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_risk_residual_evaluation_histories_created_by
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_residual_evaluation_histories;
-- +goose StatementEnd