-- +goose Up
-- +goose StatementBegin
ALTER TABLE meeting_minutes
ADD meeting_minute_status_id INT NULL;

ALTER TABLE meeting_minutes
ADD CONSTRAINT FK_meeting_minutes_meeting_minute_statuses FOREIGN KEY (meeting_minute_status_id) REFERENCES meeting_minute_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE meeting_minutes
ADD meeting_id INT NULL;

ALTER TABLE meeting_minutes
ADD CONSTRAINT FK_meeting_minutes_meetings FOREIGN KEY (meeting_id) REFERENCES meeting_minutes (id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE meeting_minutes
DROP CONSTRAINT FK_meeting_minutes_meeting_minute_statuses;

ALTER TABLE meeting_minutes
DROP COLUMN meeting_minute_status_id;

ALTER TABLE meeting_minutes
DROP CONSTRAINT FK_meeting_minutes_meetings;

ALTER TABLE meeting_minutes
DROP COLUMN meeting_id;
-- +goose StatementEnd
