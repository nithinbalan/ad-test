-- +goose Up
-- +goose StatementBegin
CREATE TABLE workflows (
    id INT PRIMARY KEY IDENTITY(1,1),
    type VARCHAR(100) NOT NULL,
     CONSTRAINT CHK_workflows_type CHECK (type IN ('document')),
    document_type_id INT NULL,
    status VARCHAR(20) NOT NULL,
     CONSTRAINT CHK_workflows_status CHECK (status IN ('active', 'inactive')), 
    created_by_id INT NOT NULL,
    updated_by_id INT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,

    CONSTRAINT FK_workflows_document_types FOREIGN KEY (document_type_id) REFERENCES document_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_workflows_user_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,  
    CONSTRAINT FK_workflows_user_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.workflows;
-- +goose StatementEnd
