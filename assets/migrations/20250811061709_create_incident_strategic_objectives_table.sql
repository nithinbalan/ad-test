-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_strategic_objectives  (
    id INT PRIMARY KEY IDENTITY(1,1),
    strategic_objective_id INT NOT NULL,
    incident_register_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    CONSTRAINT FK_incident_strategic_objectives_strategic_objectives FOREIGN KEY (strategic_objective_id) REFERENCES strategic_objectives(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_strategic_objectives_incident_register FOREIGN KEY (incident_register_id) REFERENCES incident_register(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_strategic_objectives_user_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_strategic_objectives;
-- +goose StatementEnd
