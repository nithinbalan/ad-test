-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_register_update_plan_team_members (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_register_update_plan_id INT NOT NULL,
    user_id INT NOT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    CONSTRAINT FK_risk_register_update_plan_team_members_risk_register_update_plans FOREIGN KEY (risk_register_update_plan_id) REFERENCES risk_register_update_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_register_update_plan_team_members_users_1 FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_register_update_plan_team_members_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_register_update_plan_team_members_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_register_update_plan_team_members;
-- +goose StatementEnd