-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_response_team_members (
    id INT PRIMARY KEY IDENTITY(1, 1),
    incident_response_team_id INT NOT NULL,
    user_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_incident_response_team_members_team
        FOREIGN KEY (incident_response_team_id) REFERENCES incident_response_teams(id)
        ON DELETE NO ACTION ON UPDATE CASCADE,

    CONSTRAINT FK_incident_response_team_members_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE CASCADE,

    CONSTRAINT FK_incident_response_team_members_created_by
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_response_team_members;
-- +goose StatementEnd
