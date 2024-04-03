-- +goose Up
-- +goose StatementBegin
CREATE TABLE potential_risks (
    id INT PRIMARY KEY IDENTITY(1, 1),
    risk_category_id INT NOT NULL,
    severity_level_id INT NOT NULL,
    potential_risk_status_id INT NOT NULL,
    title NVARCHAR(100) NOT NULL,
    description NVARCHAR(MAX) NULL,
    risk_source VARCHAR(45) NOT NULL,
    created_at DATETIME NOT NULL ,
    updated_at DATETIME NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    CONSTRAINT FK_potential_risks_risk_categories FOREIGN KEY (risk_category_id) REFERENCES risk_categories(id),
    CONSTRAINT FK_potential_risks_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels(id),
    CONSTRAINT FK_potential_risks_potential_risk_statuses FOREIGN KEY (potential_risk_status_id) REFERENCES potential_risk_statuses(id),
    CONSTRAINT FK_potential_risks_users1 FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT FK_potential_risks_users2 FOREIGN KEY (updated_by) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE potential_risks;
-- +goose StatementEnd
