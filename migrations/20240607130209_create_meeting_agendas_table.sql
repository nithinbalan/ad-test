-- +goose Up
-- +goose StatementBegin
CREATE TABLE meeting_agendas
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    meeting_id INT NOT NULL,
    time_from NVARCHAR(45) NULL,
    time_to NVARCHAR(45) NULL,
    title NVARCHAR(500) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    CONSTRAINT FK_meeting_agendas_meetings FOREIGN KEY (meeting_id) REFERENCES meetings (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_agendas_users_1 FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_agendas_users_2 FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE meeting_agendas;
-- +goose StatementEnd
