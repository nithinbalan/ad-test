-- +goose Up
-- +goose StatementBegin
CREATE TABLE tprm_assessment_workflow_participants (
    id INT PRIMARY KEY IDENTITY(1,1),
    tprm_assessment_workflow_id INT NOT NULL,
    user_id INT NOT NULL,
    role_id INT NULL,
    delegatee_user_id INT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_tprm_assessment_workflow_participants_tprm_workflows FOREIGN KEY (tprm_assessment_workflow_id) REFERENCES tprm_assessment_workflows(id),
    CONSTRAINT FK_tprm_assessment_workflow_participants_users_1 FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT FK_tprm_assessment_workflow_participants_users_4 FOREIGN KEY (delegatee_user_id) REFERENCES users(id),
    CONSTRAINT FK_tprm_assessment_workflow_participants_roles FOREIGN KEY (role_id) REFERENCES roles(id),
    CONSTRAINT FK_tprm_assessment_workflow_participants_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tprm_assessment_workflow_participants_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.tprm_assessment_workflow_participants;
-- +goose StatementEnd