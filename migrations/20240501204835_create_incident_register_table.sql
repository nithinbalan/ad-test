-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_register
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    reference_number VARCHAR(45) NOT NULL UNIQUE,
    incident_category_id INT NOT NULL,
    incident_sub_category_id INT NOT NULL,
    incident_type_id INT NOT NULL,
    severity_level_id INT NOT NULL,
    incident_status_id INT NOT NULL,
    title NVARCHAR(255) NOT NULL,
    description NVARCHAR(MAX) NULL,
    incident_at DATE NOT NULL,
    location VARCHAR(100),
    reported_by INT NOT NULL,
    reported_at DATE,
    quick_action NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_incident_register_incident_categories FOREIGN KEY (incident_category_id) REFERENCES incident_categories (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_register_incident_sub_categories FOREIGN KEY (incident_sub_category_id) REFERENCES incident_sub_categories (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_register_incident_types FOREIGN KEY (incident_type_id) REFERENCES incident_types (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_register_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_register_incident_statuses FOREIGN KEY (incident_status_id) REFERENCES incident_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_register_users_1 FOREIGN KEY (reported_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_register_users_2 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_register_users_3 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_register;

-- +goose StatementEnd