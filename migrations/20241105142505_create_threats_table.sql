-- +goose Up
-- +goose StatementBegin
CREATE TABLE threats
(
    id INT PRIMARY KEY IDENTITY(1,1),
    title NVARCHAR(500) NOT NULL,
    threat_type_id INT NOT NULL,
    severity_level_id INT NOT NULL,
    threat_status_id INT NOT NULL,
    responsible_user_id INT NOT NULL,
    responsible_unit_id INT NOT NULL,
    description NVARCHAR(MAX),
    detection_date DATE,
    resolution_date DATE,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT,
    updated_at DATETIME,

    CONSTRAINT FK_threats_threat_types FOREIGN KEY (threat_type_id) REFERENCES threat_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_threats_threat_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_threats_threat_statuses FOREIGN KEY (threat_status_id) REFERENCES threat_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_threats_organization_units FOREIGN KEY (responsible_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_threats_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_threats_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_threats_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS threats;
-- +goose StatementEnd
