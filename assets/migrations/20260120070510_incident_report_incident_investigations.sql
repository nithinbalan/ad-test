-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_report_incident_investigations (
    id INT PRIMARY KEY IDENTITY(1,1),

    incident_report_id INT NOT NULL,
    incident_report_incident_id INT NOT NULL,
    incident_investigation_id INT NOT NULL,

    started_at DATETIME NOT NULL,
    investigation_details NVARCHAR(MAX) NULL,
    action_recommendation NVARCHAR(MAX) NULL,
    
    team_lead_name_en NVARCHAR(255) NOT NULL,
    team_lead_name_ar NVARCHAR(255) NOT NULL,

    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_incident_report_inv_reports FOREIGN KEY (incident_report_id) REFERENCES incident_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_report_inv_incidents FOREIGN KEY (incident_report_incident_id) REFERENCES incident_report_incidents (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_report_inv_investigations FOREIGN KEY (incident_investigation_id) REFERENCES incident_investigations (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_report_incident_investigations;
-- +goose StatementEnd