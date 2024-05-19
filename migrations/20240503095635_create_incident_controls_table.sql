-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_controls (
    id INT PRIMARY KEY IDENTITY(1, 1),
    incident_register_id INT NOT NULL,
    control_id INT NOT NULL,
    created_by INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_incident_controls_incident_register FOREIGN KEY (incident_register_id) REFERENCES incident_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_controls_controls FOREIGN KEY (control_id) REFERENCES controls(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_controls;

-- +goose StatementEnd