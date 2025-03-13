-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_action_plan_controls (
    id INT PRIMARY KEY IDENTITY(1, 1),
    risk_action_plan_id INT NOT NULL,
    control_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_risk_action_plan_controls_risk_action_plans FOREIGN KEY (risk_action_plan_id) REFERENCES risk_action_plans (id) ON DELETE NO ACTION ON UPDATE NO ACTION ,
    CONSTRAINT FK_risk_action_plan_controls_controls FOREIGN KEY (control_id) REFERENCES controls (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_action_plan_controls_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_action_plan_controls_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_action_plan_controls;
-- +goose StatementEnd
