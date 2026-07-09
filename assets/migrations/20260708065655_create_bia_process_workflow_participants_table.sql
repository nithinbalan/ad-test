-- +goose Up
-- +goose StatementBegin
CREATE TABLE bia_process_workflow_participants (
    id INT PRIMARY KEY IDENTITY(1,1),
    bia_process_workflow_id INT NOT NULL,
    user_id INT NOT NULL,
    role_id INT NULL,
    delegatee_user_id INT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_bia_process_wp_workflow
        FOREIGN KEY (bia_process_workflow_id)
        REFERENCES bia_process_workflows(id),

    CONSTRAINT FK_bia_process_wp_user
        FOREIGN KEY (user_id)
        REFERENCES users(id),

    CONSTRAINT FK_bia_process_wp_delegatee_user
        FOREIGN KEY (delegatee_user_id)
        REFERENCES users(id),

    CONSTRAINT FK_bia_process_wp_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id),

    CONSTRAINT FK_bia_process_wp_users_1
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_bia_process_wp_users_2
        FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE bia_process_workflow_participants;
-- +goose StatementEnd
