-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_closure_workflow_histories (
    id INT PRIMARY KEY IDENTITY(1,1),
    incident_closure_id INT NOT NULL,
    level INT NOT NULL,
    workflow_status_id INT NOT NULL,
    comment NVARCHAR(MAX),
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,

    CONSTRAINT FK_incident_closure_workflow_histories_workflow
        FOREIGN KEY (incident_closure_id)
        REFERENCES incident_closures(id)
       ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_incident_closure_workflow_histories_status
        FOREIGN KEY (workflow_status_id)
        REFERENCES workflow_statuses(id)
       ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_incident_closure_workflow_histories_users_1
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
       ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_incident_closure_workflow_histories_users_2
        FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
       ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_closure_workflow_histories;
-- +goose StatementEnd