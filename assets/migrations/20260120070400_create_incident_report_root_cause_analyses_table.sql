-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_report_root_cause_analyses (
    id INT PRIMARY KEY IDENTITY(1,1),

    incident_report_incident_id INT NOT NULL,
    incident_report_id INT NOT NULL,

    why INT NOT NULL,
    question NVARCHAR(500) NOT NULL,
    answer NVARCHAR(1000) NULL,

    CONSTRAINT FK_incident_report_rca_incident FOREIGN KEY (incident_report_incident_id) REFERENCES incident_report_incidents (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_report_rca_report FOREIGN KEY (incident_report_id) REFERENCES incident_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_report_root_cause_analyses;
-- +goose StatementEnd
