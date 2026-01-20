-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_report_sections (
    id INT PRIMARY KEY IDENTITY(1,1),

    incident_report_id INT NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('executive-summary', 'scope', 'objectives', 'overview', 'details', 'investigation', 'incident-closure', 'post-incident-review','rca','communication-log','action-plan')),

    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,

    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,

    order_num INT NOT NULL,

    CONSTRAINT FK_incident_report_sections_incident_reports FOREIGN KEY (incident_report_id) REFERENCES incident_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_report_sections;
-- +goose StatementEnd
