-- +goose Up
-- +goose StatementBegin
ALTER TABLE projects
ADD
    workflow_status_id INT NULL,
    submitted_by_id INT NULL,
    submitted_at DATETIME NULL,
    next_review_level INT NULL,
    approved_by_id INT NULL,
    approved_at DATETIME NULL,
    rejected_by_id INT NULL,
    rejected_at DATETIME NULL,
    incident_action_plan_id INT NULL,

    CONSTRAINT FK_projects_workflow_statuses
        FOREIGN KEY (workflow_status_id) REFERENCES workflow_statuses(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_projects_users_4
        FOREIGN KEY (submitted_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_projects_users_5
        FOREIGN KEY (approved_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_projects_users_6
        FOREIGN KEY (rejected_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_projects_incident_action_plans
        FOREIGN KEY (incident_action_plan_id) REFERENCES incident_action_plans(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE projects
DROP CONSTRAINT FK_projects_workflow_statuses;

ALTER TABLE projects
DROP CONSTRAINT FK_projects_users_4;

ALTER TABLE projects
DROP CONSTRAINT FK_projects_users_5;

ALTER TABLE projects
DROP CONSTRAINT FK_projects_users_6;

ALTER TABLE projects
DROP CONSTRAINT FK_projects_incident_action_plans;

ALTER TABLE projects
DROP COLUMN workflow_status_id,
            submitted_by_id,
            submitted_at,
            next_review_level,
            approved_by_id,
            approved_at,
            rejected_by_id,
            rejected_at,
            incident_action_plan_id;
-- +goose StatementEnd
