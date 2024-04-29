-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_matrix_impact_guidelines
(
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_impact_category_id INT NOT NULL,
    risk_matrix_impact_id INT NOT NULL,
    guideline VARCHAR(500) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_risk_matrix_impact_guidelines_status CHECK (status IN ('active', 'inactive')),
    created_by INT,
    updated_by INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_risk_matrix_impact_guidelines_risk_impact_categories FOREIGN KEY (risk_impact_category_id) REFERENCES risk_impact_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_matrix_impact_guidelines_risk_matrix_impacts FOREIGN KEY (risk_matrix_impact_id) REFERENCES risk_matrix_impacts(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_matrix_impact_guidelines_users_1 FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_matrix_impact_guidelines_users_2 FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_matrix_impact_guidelines;
-- +goose StatementEnd
