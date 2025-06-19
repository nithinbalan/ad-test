-- +goose Up
-- +goose StatementBegin
CREATE TABLE meeting_agendas
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    meeting_id INT NOT NULL,
    title NVARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_meeting_agendas_meetings FOREIGN KEY (meeting_id) REFERENCES meetings (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_agendas_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_agendas_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE meeting_agendas;
-- +goose StatementEnd
