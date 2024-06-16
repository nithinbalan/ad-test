-- +goose Up
-- +goose StatementBegin
CREATE TABLE sla_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    sla_register_id INT NOT NULL,
    version VARCHAR(100) NOT NULL,
    issue_date DATE NULL,
    expiry_date DATE NULL,
    comment VARCHAR(MAX) NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    is_latest TINYINT DEFAULT 1,
    sa1_date DATE DEFAULT NULL,
    sa2_date DATE DEFAULT NULL,
    CONSTRAINT FK_sla_documents_sla_register FOREIGN KEY (sla_register_id) REFERENCES sla_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_documents_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_documents_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.sla_documents;
-- +goose StatementEnd
