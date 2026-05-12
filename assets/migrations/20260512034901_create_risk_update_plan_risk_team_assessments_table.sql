-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_update_plan_risk_team_assessments(
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_update_plan_timeline_id INT NOT NULL,
    risk_champion_id INT NOT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    CONSTRAINT FK_risk_update_plan_risk_team_assessments_risk_update_plan_timelines FOREIGN KEY (risk_update_plan_timeline_id) REFERENCES risk_update_plan_timelines(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_update_plan_risk_team_assessments_users_1 FOREIGN KEY (risk_champion_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_update_plan_risk_team_assessments_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_update_plan_risk_team_assessments_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_update_plan_risk_team_assessments;
-- +goose StatementEnd
