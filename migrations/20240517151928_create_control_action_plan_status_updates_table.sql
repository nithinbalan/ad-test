-- +goose Up
-- +goose StatementBegin
CREATE TABLE control_action_plan_status_updates (
    id INT PRIMARY KEY IDENTITY(1, 1),
    control_action_plan_id INT NOT NULL,
    action_plan_status_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    percentage INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    created_by INT NOT NULL,
    updated_by INT,
    CONSTRAINT FK_control_action_plan_status_updates_control_action_plans FOREIGN KEY (control_action_plan_id) REFERENCES control_action_plans (id) ON DELETE NO ACTION ON UPDATE NO ACTION ,
    CONSTRAINT FK_control_action_plan_status_updates_action_plan_statuses FOREIGN KEY (action_plan_status_id) REFERENCES action_plan_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_control_action_plan_status_updates_users_1 FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_control_action_plan_status_updates_users_2 FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE control_action_plan_status_updates;
-- +goose StatementEnd
