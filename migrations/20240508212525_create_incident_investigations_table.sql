-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_investigations
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    incident_register_id INT NOT NULL,
    initiated_at DATETIME NOT NULL,
    description NVARCHAR(MAX) NULL,
    team_lead_id INT NOT NULL,
    action_recommendation NVARCHAR(MAX) NULL,
    incident_investigation_status_id INT NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_incident_investigations_incident_register FOREIGN KEY (incident_register_id) REFERENCES incident_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_investigations_incident_investigation_statuses FOREIGN KEY (incident_investigation_status_id) REFERENCES incident_investigation_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_investigations_users_1 FOREIGN KEY (team_lead_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_investigations_users_2 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_investigations_users_3 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_investigations;

-- +goose StatementEnd