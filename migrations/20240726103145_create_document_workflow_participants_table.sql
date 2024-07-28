-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_workflow_participants (
    id INT PRIMARY KEY IDENTITY(1,1),
    document_workflow_id INT NOT NULL,
    workflow_level_type_id INT NOT NULL,
    user_id INT NOT NULL,
    role_id INT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_document_workflow_participants_document_workflows FOREIGN KEY (document_workflow_id) REFERENCES document_workflows(id),
    CONSTRAINT FK_document_workflow_participants_workflow_level_types FOREIGN KEY (workflow_level_type_id) REFERENCES workflow_level_types(id),
    CONSTRAINT FK_document_workflow_participants_users_1 FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT FK_document_workflow_participants_roles FOREIGN KEY (role_id) REFERENCES roles(id),
    CONSTRAINT FK_document_workflow_participants_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_workflow_participants_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_workflow_participants;
-- +goose StatementEnd
