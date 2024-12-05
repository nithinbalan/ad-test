-- +goose Up
-- +goose StatementBegin
CREATE TABLE vulnerabilities (
    id INT PRIMARY KEY IDENTITY(1,1),
    reference_number VARCHAR(45) NOT NULL UNIQUE,
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX),
    vulnerability_type_id INT NOT NULL,
    severity_level_id INT,
    vulnerability_status_id INT NOT NULL,
    detection_date DATE NOT NULL,
    resolution_date DATE,
    responsible_user_id INT NOT NULL,
    responsible_unit_id INT NOT NULL,
    source NVARCHAR(500),
    created_by_id INT NOT NULL,
    updated_by_id INT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,

    CONSTRAINT FK_vulnerabilities_vulnerability_types FOREIGN KEY (vulnerability_type_id) REFERENCES vulnerability_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_vulnerabilities_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_vulnerabilities_vulnerability_statuses FOREIGN KEY (vulnerability_status_id) REFERENCES vulnerability_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_vulnerabilities_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_vulnerabilities_organization_units FOREIGN KEY (responsible_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_vulnerabilities_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_vulnerabilities_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS vulnerabilities;
-- +goose StatementEnd
