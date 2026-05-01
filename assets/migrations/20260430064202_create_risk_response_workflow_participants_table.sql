-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_response_workflow_participants (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_response_workflow_id INT NOT NULL,
    user_id INT NOT NULL,
    role_id INT NULL,
    delegatee_user_id INT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_rrwp_workflow
        FOREIGN KEY (risk_response_workflow_id)
        REFERENCES risk_response_workflows(id),

    CONSTRAINT FK_rrwp_user
        FOREIGN KEY (user_id)
        REFERENCES users(id),

    CONSTRAINT FK_rrwp_delegatee_user
        FOREIGN KEY (delegatee_user_id)
        REFERENCES users(id),

    CONSTRAINT FK_rrwp_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id),

    CONSTRAINT FK_rrwp_users_1
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_rrwp_users_2
        FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.risk_response_workflow_participants;
-- +goose StatementEnd
