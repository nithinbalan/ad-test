-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_register_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,
    icon VARCHAR(10) NOT NULL,
    created_by INT,
    updated_by INT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_risk_documents_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_documents_users_1 FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_documents_users_2 FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.risk_documents;
-- +goose StatementEnd
