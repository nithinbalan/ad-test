-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_timelines (
    id INT PRIMARY KEY IDENTITY(1, 1),
    incident_register_id INT NOT NULL,
    incident_timeline_type_id INT NOT NULL,
    timeline_status_id INT NOT NULL,
    performed_by_id INT NOT NULL,
    performed_at DATETIME NULL,
    CONSTRAINT FK_incident_timelines_incident_register FOREIGN KEY (incident_register_id) REFERENCES incident_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_timelines_incident_timeline_types FOREIGN KEY (incident_timeline_type_id) REFERENCES incident_timeline_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_timelines_timeline_statuses FOREIGN KEY (timeline_status_id) REFERENCES timeline_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_timelines_users_1 FOREIGN KEY (performed_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION

);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_timelines;

-- +goose StatementEnd
