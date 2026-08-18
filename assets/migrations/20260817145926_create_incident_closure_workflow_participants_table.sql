-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_closure_workflow_participants (
    id INT PRIMARY KEY IDENTITY(1,1),
    incident_closure_workflow_id INT NOT NULL,
    user_id INT NOT NULL,
    role_id INT NULL,
    delegatee_user_id INT NULL,
    workflow_status_id INT NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_incident_closure_workflow_participants_incident_closure_workflows FOREIGN KEY (incident_closure_workflow_id) REFERENCES incident_closure_workflows(id),
    CONSTRAINT FK_incident_closure_workflow_participants_users_1 FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT FK_incident_closure_workflow_participants_users_4 FOREIGN KEY (delegatee_user_id) REFERENCES users(id),
    CONSTRAINT FK_incident_closure_workflow_participants_roles FOREIGN KEY (role_id) REFERENCES roles(id),
    CONSTRAINT FK_incident_closure_workflow_participants_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_closure_workflow_participants_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_closure_workflow_participants_workflow_statuses FOREIGN KEY (workflow_status_id) REFERENCES workflow_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_closure_workflow_participants;
-- +goose StatementEnd
