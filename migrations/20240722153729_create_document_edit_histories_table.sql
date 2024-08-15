-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_edit_histories (
    id INT PRIMARY KEY IDENTITY(1,1),
    document_id INT NOT NULL,
    document_version_id INT NOT NULL,
    action_type VARCHAR(20) NOT NULL DEFAULT 'create',
        CONSTRAINT CHK_document_edit_histories_action_type CHECK (action_type IN ('create', 'update','check-in', 'check-out', 'revert-checkout')),
    comment NVARCHAR(MAX) NULL,    
    url VARCHAR(1000) NULL,
    name NVARCHAR(255) NULL,
    size INT NULL,
    token VARCHAR(255) NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_document_edit_histories_documents FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_edit_histories_document_versions FOREIGN KEY (document_version_id) REFERENCES document_versions(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_edit_histories_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_edit_histories;
-- +goose StatementEnd
