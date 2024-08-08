-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_controls (
    id INT PRIMARY KEY IDENTITY(1,1),
    document_id INT NOT NULL,
    document_version_id INT NULL,
    document_version_clause_id INT NULL,
    control_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_document_controls_document FOREIGN KEY (document_id) REFERENCES documents(id),
    CONSTRAINT FK_document_controls_document_version FOREIGN KEY (document_version_id) REFERENCES document_versions(id),
    CONSTRAINT FK_document_controls_document_version_clause FOREIGN KEY (document_version_clause_id) REFERENCES document_version_clauses(id),
    CONSTRAINT FK_document_controls_controls FOREIGN KEY (control_id) REFERENCES controls(id),
    CONSTRAINT FK_document_controls_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id),
    CONSTRAINT FK_document_controls_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id)  
    
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE document_controls;
-- +goose StatementEnd
