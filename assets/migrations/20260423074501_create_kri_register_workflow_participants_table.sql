-- +goose Up
-- +goose StatementBegin
CREATE TABLE kri_register_workflow_participants (
    id INT PRIMARY KEY IDENTITY(1,1),
    kri_register_workflow_id INT NOT NULL,
    user_id INT NOT NULL,
    role_id INT NULL,
    delegatee_user_id INT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_krwp_workflow
        FOREIGN KEY (kri_register_workflow_id)
        REFERENCES kri_register_workflows(id),

    CONSTRAINT FK_krwp_user
        FOREIGN KEY (user_id)
        REFERENCES users(id),

    CONSTRAINT FK_krwp_delegatee_user
        FOREIGN KEY (delegatee_user_id)
        REFERENCES users(id),

    CONSTRAINT FK_krwp_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.kri_register_workflow_participants;
-- +goose StatementEnd