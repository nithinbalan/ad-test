-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_action_plan_status_updates (
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_action_plan_id INT NOT NULL,
    action_plan_status_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    percentage INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    CONSTRAINT FK_compliance_action_plan_status_updates_compliance_action_plans FOREIGN KEY (compliance_action_plan_id) REFERENCES compliance_action_plans (id),
    CONSTRAINT FK_compliance_action_plan_status_updates_action_plan_statuses FOREIGN KEY (action_plan_status_id) REFERENCES action_plan_statuses (id),
    CONSTRAINT FK_compliance_action_plan_status_updates_users1 FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT FK_compliance_action_plan_status_updates_users2 FOREIGN KEY (updated_by) REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_action_plan_status_updates;
-- +goose StatementEnd
