-- +goose Up
-- +goose StatementBegin

CREATE TABLE meeting_participants (
    id INT PRIMARY KEY IDENTITY(1, 1),
    meeting_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    CONSTRAINT FK_meeting_participants_meetings FOREIGN KEY (meeting_id) REFERENCES meetings (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_participants_users_1 FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_participants_users_2 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS meeting_participants;

-- +goose StatementEnd
