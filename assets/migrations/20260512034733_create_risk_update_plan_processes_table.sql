-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_update_plan_processes (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_update_plan_timeline_id INT NOT NULL,
    process_id INT NOT NULL,
    created_by_id INT,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_risk_update_plan_processes_risk_update_plan_timelines FOREIGN KEY (risk_update_plan_timeline_id) REFERENCES risk_update_plan_timelines(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_update_plan_processes_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_update_plan_processes_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_update_plan_processes;
-- +goose StatementEnd
