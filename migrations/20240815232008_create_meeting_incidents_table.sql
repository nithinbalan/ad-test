-- +goose Up
-- +goose StatementBegin
CREATE TABLE meeting_incidents (
    id INT PRIMARY KEY IDENTITY(1, 1),
    meeting_id INT NOT NULL,
    incident_register_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_meeting_incident_meetings FOREIGN KEY (meeting_id) REFERENCES meetings(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_incident_incident_register FOREIGN KEY (incident_register_id) REFERENCES incident_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_incident_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE meeting_incidents;

-- +goose StatementEnd