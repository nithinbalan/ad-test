-- +goose Up
-- +goose StatementBegin
CREATE TABLE meeting_action_plans (
    id INT IDENTITY(1,1) PRIMARY KEY,
    meeting_id INT NOT NULL,
    meeting_agenda_id INT NULL,
    title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX),
    severity_level_id INT,
    target_date DATE,
    responsible_user_id INT NOT NULL,
    action_plan_status_id INT NOT NULL,
    percentage INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    created_by_id INT NOT NULL,
    updated_by_id INT,
    CONSTRAINT FK_meeting_action_plans_meetings FOREIGN KEY (meeting_id) REFERENCES meetings(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_action_plans_meeting_agendas FOREIGN KEY (meeting_id) REFERENCES meeting_agendas(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_action_plans_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_action_plans_user_1 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_action_plans_action_plan_statuses FOREIGN KEY (action_plan_status_id) REFERENCES action_plan_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_action_plans_users_2 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_action_plans_users_3 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE meeting_action_plans;
-- +goose StatementEnd
