-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_action_plans (
    id INT PRIMARY KEY IDENTITY(1, 1),
    project_milestone_id INT NOT NULL,
    title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX) NULL,
    target_date DATE NULL,
    responsible_user_id INT NOT NULL,
    percentage INT NOT NULL DEFAULT 0, 
    comments NVARCHAR(MAX) NULL,
    action_plan_status_id INT NOT NULL,
    created_at DATETIME NOT NULL ,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_project_action_plans_project_milestones FOREIGN KEY (project_milestone_id) REFERENCES project_milestones (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_project_action_plans_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_project_action_plans_action_plan_statuses FOREIGN KEY (action_plan_status_id) REFERENCES action_plan_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_project_action_plans_users_2 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_project_action_plans_users_3 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_action_plans;
-- +goose StatementEnd
