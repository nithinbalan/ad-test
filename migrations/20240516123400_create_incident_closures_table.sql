-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_closures
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    incident_register_id INT NOT NULL,
    completion_date DATE NOT NULL,
    resolution_summary NVARCHAR(MAX) NULL,
    notes NVARCHAR(MAX) NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_incident_closures_incident_register FOREIGN KEY (incident_register_id) REFERENCES incident_register (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_closures_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_closures_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.incident_closures;
-- +goose StatementEnd
