-- +goose Up
-- +goose StatementBegin

CREATE TABLE meetings (
    id INT PRIMARY KEY IDENTITY(1, 1),
    meeting_type_id INT NOT NULL,
    meeting_frequency_id INT NULL,
    meeting_id INT NULL,
    title VARCHAR(100) NOT NULL,
    location VARCHAR(100) NULL,
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    organizer_id INT NOT NULL,
    created_by INT NOT NULL,
    updated_by INT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    meeting_status_id INT NOT NULL,
    CONSTRAINT FK_meetings_meeting_types FOREIGN KEY (meeting_type_id) REFERENCES meeting_types (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meetings_self FOREIGN KEY (meeting_id) REFERENCES meetings (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meetings_meeting_frequencies FOREIGN KEY (meeting_frequency_id) REFERENCES meeting_frequencies (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meetings_users_1 FOREIGN KEY (organizer_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meetings_users_2 FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meetings_users_3 FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meetings_meeting_statuses FOREIGN KEY (meeting_status_id) REFERENCES meeting_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE dbo.meetings;

-- +goose StatementEnd
