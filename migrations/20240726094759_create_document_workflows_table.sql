-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_workflows (
    id INT PRIMARY KEY IDENTITY(1,1),
    document_version_id INT NOT NULL,
    level INT NOT NULL,
    approval_types VARCHAR(255) NOT NULL DEFAULT 'anyone',
      CONSTRAINT CHK_document_workflows_approval_types CHECK (approval_types IN ('anyone', 'all')),
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_document_workflows_document_versions FOREIGN KEY (document_version_id) REFERENCES document_versions(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_workflows_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_workflows_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_workflows;
-- +goose StatementEnd
