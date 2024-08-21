-- +goose Up
-- +goose StatementBegin
ALTER TABLE meeting_action_plans
ADD comments NVARCHAR(MAX) NULL;


ALTER TABLE meeting_action_plans
ADD meeting_agenda_id INT NULL;

ALTER TABLE meeting_action_plans
ADD CONSTRAINT FK_meeting_action_plans_meeting_agendas FOREIGN KEY (meeting_agenda_id) REFERENCES meeting_agendas(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE meeting_action_plans
DROP CONSTRAINT FK_meeting_action_plans_meeting_agendas;

ALTER TABLE meeting_action_plans
DROP COLUMN meeting_agenda_id;
-- +goose StatementEnd
