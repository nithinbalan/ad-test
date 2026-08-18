-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_register_response_teams (
    id INT PRIMARY KEY IDENTITY(1, 1),
    incident_register_id INT NOT NULL,
    user_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_incident_register_response_teams_incident_register
        FOREIGN KEY (incident_register_id) REFERENCES incident_register(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_incident_register_response_teams_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_incident_register_response_teams_created_by
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_register_response_teams;
-- +goose StatementEnd
