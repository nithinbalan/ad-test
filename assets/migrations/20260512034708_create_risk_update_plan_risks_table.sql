-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_update_plan_risks (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_update_plan_timeline_id INT NOT NULL,
    risk_register_id INT NOT NULL,
    created_by_id INT,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_risk_update_plan_risks_risk_update_plan_timelines FOREIGN KEY (risk_update_plan_timeline_id) REFERENCES risk_update_plan_timelines(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_update_plan_risks_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_update_plan_risks_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_update_plan_risks;
-- +goose StatementEnd
