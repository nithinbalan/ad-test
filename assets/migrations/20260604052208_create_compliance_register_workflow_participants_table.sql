-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_register_workflow_participants (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_register_workflow_id INT NOT NULL,
    user_id INT NOT NULL,
    role_id INT NULL,
    delegatee_user_id INT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_compliance_register_workflow_participants_workflow
        FOREIGN KEY (compliance_register_workflow_id)
        REFERENCES compliance_register_workflows(id),

    CONSTRAINT FK_compliance_register_workflow_participants_user
        FOREIGN KEY (user_id)
        REFERENCES users(id),

    CONSTRAINT FK_compliance_register_workflow_participants_delegatee_user
        FOREIGN KEY (delegatee_user_id)
        REFERENCES users(id),

    CONSTRAINT FK_compliance_register_workflow_participants_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_register_workflow_participants;
-- +goose StatementEnd