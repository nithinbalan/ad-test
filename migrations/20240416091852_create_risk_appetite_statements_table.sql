-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_appetite_statements (
    id INT PRIMARY KEY IDENTITY(1, 1),
    reference_number VARCHAR(45) NULL,
    risk_category_id INT NOT NULL,
    statement_en NVARCHAR(MAX) NOT NULL,
    statement_ar NVARCHAR(MAX) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_risk_appetite_statements_status CHECK (status IN ('active', 'inactive')),
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    CONSTRAINT FK_risk_appetite_statements_risk_categories FOREIGN KEY (risk_category_id) REFERENCES risk_categories (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_appetite_statements_users_1 FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_appetite_statements_users_2 FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS risk_appetite_statements;
-- +goose StatementEnd
