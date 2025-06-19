-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_timelines (
    id INT PRIMARY KEY IDENTITY(1, 1),
    risk_register_id INT NOT NULL,
    risk_timeline_type_id INT NOT NULL,
    timeline_status_id INT NOT NULL,
    performed_by_id INT NULL,
    performed_at DATETIME NULL,
    CONSTRAINT FK_risk_timelines_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_timelines_risk_timeline_types FOREIGN KEY (risk_timeline_type_id) REFERENCES risk_timeline_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_timelines_timeline_statuses FOREIGN KEY (timeline_status_id) REFERENCES timeline_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_timelines_users_1 FOREIGN KEY (performed_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION

);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_timelines;

-- +goose StatementEnd
