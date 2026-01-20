-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_report_incidents (
    id INT PRIMARY KEY IDENTITY(1,1),

    incident_report_id INT NOT NULL,
    incident_register_id INT NOT NULL,
    reference_number NVARCHAR(50) NULL, 

    incident_title NVARCHAR(225) NOT NULL,
    incident_description NVARCHAR(225) NULL,
    incident_details NVARCHAR(MAX) NULL,

    incident_category_title_en NVARCHAR(255) NULL,
    incident_category_title_ar NVARCHAR(255) NULL,
    incident_sub_category_title_en NVARCHAR(255) NULL,
    incident_sub_category_title_ar NVARCHAR(255) NULL,

    severity_level_title_en NVARCHAR(255) NULL,
    severity_level_title_ar NVARCHAR(255) NULL,
    severity_level_color_code NVARCHAR(255) NULL,

    responsible_user_name_en NVARCHAR(255) NOT NULL,
    responsible_user_name_ar NVARCHAR(255) NOT NULL,
    reported_user_name_en NVARCHAR(255) NOT NULL,
    reported_user_name_ar NVARCHAR(255) NOT NULL,

    incident_at DATETIME NOT NULL,
    reported_at DATETIME NOT NULL,
    incident_location NVARCHAR(255) NOT NULL,

    incident_status_title_en NVARCHAR(255) NULL,
    incident_status_title_ar NVARCHAR(255) NULL,

    created_at DATETIME NULL,
    updated_at DATETIME NULL,
    closed_date DATETIME NULL,

    CONSTRAINT FK_incident_report_incidents_incident_reports FOREIGN KEY (incident_report_id) REFERENCES incident_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_report_incidents_incident_register FOREIGN KEY (incident_register_id) REFERENCES incident_register (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_report_incidents;
-- +goose StatementEnd
