-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_change_request_workflow_participants (
    id INT PRIMARY KEY IDENTITY(1,1),
    document_change_request_workflow_id INT NOT NULL,
    user_id INT NOT NULL,
    role_id INT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_document_change_request_workflow_participants_workflows FOREIGN KEY (document_change_request_workflow_id) REFERENCES document_change_request_workflows(id),
    CONSTRAINT FK_document_change_request_workflow_participants_users_1 FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT FK_document_change_request_workflow_participants_roles FOREIGN KEY (role_id) REFERENCES roles(id),
    CONSTRAINT FK_document_change_request_workflow_participants_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_document_change_request_workflow_participants_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_change_request_workflow_participants;
-- +goose StatementEnd
