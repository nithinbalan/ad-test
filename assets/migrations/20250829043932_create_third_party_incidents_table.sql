-- +goose Up
-- +goose StatementBegin
CREATE TABLE third_party_incidents
(
    id INT PRIMARY KEY IDENTITY(1, 1),
	third_party_register_id INT NOT NULL,
    incident_register_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
	CONSTRAINT FK_third_party_incidents_third_party_register FOREIGN KEY (third_party_register_id) REFERENCES third_party_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_incidents_incident_register FOREIGN KEY (incident_register_id) REFERENCES incident_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
	CONSTRAINT FK_third_party_incidents_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
   
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE third_party_incidents;

-- +goose StatementEnd