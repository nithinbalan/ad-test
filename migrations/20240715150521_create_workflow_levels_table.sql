-- +goose Up
-- +goose StatementBegin
CREATE TABLE workflow_levels (
    id INT PRIMARY KEY IDENTITY(1,1),
    workflow_id INT NOT NULL,
    workflow_role_type_id INT NOT NULL,
    level INT NOT NULL,
    approval_type VARCHAR(50) NOT NULL,
      CONSTRAINT CHK_workflow_levels_approval_type CHECK (approval_type IN ('anyone', 'all')),
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    
    CONSTRAINT FK_workflow_levels_workflow FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_workflow_levels_workflow_role_type FOREIGN KEY (workflow_role_type_id) REFERENCES workflow_role_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_workflow_levels_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_workflow_levels_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workflow_levels;
-- +goose StatementEnd
