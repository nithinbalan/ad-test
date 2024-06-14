-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_register (
    id INT PRIMARY KEY IDENTITY(1, 1),
    reference_number VARCHAR(45) NOT NULL UNIQUE,
    risk_category_id INT NOT NULL,
    risk_source_id INT NOT NULL,
    owner_unit_id INT NOT NULL,
    owner_id INT NOT NULL,
    risk_status_id INT NOT NULL,
    title NVARCHAR(255) NOT NULL,
    description NVARCHAR(MAX) NULL,
    reported_at DATE NULL,
    current_likelihood_id INT NULL,
    current_impact_id INT NULL,
    current_risk_rating_id INT NULL,
    current_risk_rating_value_id INT NULL,
    review_frequency_id INT NULL,
    next_review_date DATE NULL,
    comments NVARCHAR(MAX) NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_risk_register_risk_categories FOREIGN KEY (risk_category_id) REFERENCES risk_categories (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_register_risk_sources FOREIGN KEY (risk_source_id) REFERENCES risk_sources (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_register_users_1 FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_register_organization_units FOREIGN KEY (owner_unit_id) REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_register_review_frequencies FOREIGN KEY (review_frequency_id) REFERENCES review_frequencies (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_register_risk_matrix_likelihoods FOREIGN KEY (current_likelihood_id) REFERENCES risk_matrix_likelihoods(id),
    CONSTRAINT FK_risk_register_risk_matrix_impacts FOREIGN KEY (current_impact_id) REFERENCES risk_matrix_impacts(id),
    CONSTRAINT FK_risk_register_risk_matrix_risk_ratings FOREIGN KEY (current_risk_rating_id) REFERENCES risk_matrix_risk_ratings(id),
    CONSTRAINT FK_risk_register_risk_matrix_risk_rating_values FOREIGN KEY (current_risk_rating_value_id) REFERENCES risk_matrix_risk_rating_values(id),
    CONSTRAINT FK_risk_register_users_2 FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_register_users_3 FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.risk_register;
-- +goose StatementEnd
