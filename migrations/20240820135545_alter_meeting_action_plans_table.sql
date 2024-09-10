-- +goose Up
-- +goose StatementBegin
ALTER TABLE meeting_action_plans
ADD comments NVARCHAR(MAX) NULL;


ALTER TABLE meeting_action_plans
ADD meeting_minute_id INT NULL;

ALTER TABLE meeting_action_plans
ADD CONSTRAINT FK_meeting_action_plans_meeting_minutes FOREIGN KEY (meeting_minute_id) REFERENCES meeting_minutes(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE meeting_action_plans
DROP CONSTRAINT FK_meeting_action_plans_meeting_minutes;

ALTER TABLE meeting_action_plans
DROP COLUMN meeting_minute_id;
-- +goose StatementEnd
