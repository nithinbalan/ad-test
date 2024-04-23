-- +goose Up
-- +goose StatementBegin
CREATE TABLE sla_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    sla_register_id INT NOT NULL,
    version VARCHAR(100) NOT NULL,
    issue_date DATE NULL,
    expiry_date DATE NULL,
    comment VARCHAR(MAX) NULL,
    url VARCHAR(255) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    type VARCHAR(255) DEFAULT NULL,
    icon VARCHAR(10) NOT NULL, -- Assuming icon names like 'pdf', 'excel', etc.
        CONSTRAINT CHK_sla_documents_icon CHECK (icon IN ('pdf', 'excel', 'ppt', 'image', 'text', 'other')),
    created_by INT,
    updated_by INT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_latest TINYINT DEFAULT 1,
    sa1_date DATE DEFAULT NULL,
    sa2_date DATE DEFAULT NULL,
    CONSTRAINT FK_sla_documents_sla_register FOREIGN KEY (sla_register_id) REFERENCES sla_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_documents_users_1 FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_documents_users_2 FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.sla_documents;
-- +goose StatementEnd
