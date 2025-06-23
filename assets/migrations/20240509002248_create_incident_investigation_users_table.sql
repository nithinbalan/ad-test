-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_investigation_users
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    incident_investigation_id INT NOT NULL,
    user_id INT NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_incident_investigation_users_incident_investigations FOREIGN KEY (incident_investigation_id) REFERENCES incident_investigations(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_investigation_users_users_1 FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_investigations_users_users_2 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_investigations_users_users_3 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION

);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_investigation_users;

-- +goose StatementEnd