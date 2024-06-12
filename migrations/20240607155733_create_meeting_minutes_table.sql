-- +goose Up
-- +goose StatementBegin
CREATE TABLE meeting_minutes
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    meeting_agenda_id INT NULL,
    minutes NVARCHAR(MAX) NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    CONSTRAINT FK_meeting_minutes_meeting_agendas FOREIGN KEY (meeting_agenda_id) REFERENCES meeting_agendas (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_minutes_users_1 FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_minutes_users_2 FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE meeting_minutes;
-- +goose StatementEnd
