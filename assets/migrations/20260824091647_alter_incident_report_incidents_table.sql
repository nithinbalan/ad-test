-- +goose Up
-- +goose StatementBegin
ALTER TABLE incident_report_incidents
ADD description NVARCHAR(MAX) NULL,
    quick_action NVARCHAR(MAX) NULL,

    incident_type_title_en NVARCHAR(255) NULL,
    incident_type_title_ar NVARCHAR(255) NULL,

    incident_status_color_code NVARCHAR(255) NULL,

    responsible_unit_title_en NVARCHAR(255) NULL,
    responsible_unit_title_ar NVARCHAR(255) NULL,
    responsible_unit_color_code NVARCHAR(255) NULL,

    business_owner_name_en NVARCHAR(255) NULL,
    business_owner_name_ar NVARCHAR(255) NULL,

    response_team_lead_name_en NVARCHAR(255) NULL,
    response_team_lead_name_ar NVARCHAR(255) NULL,
    response_team_members_en NVARCHAR(MAX) NULL,
    response_team_members_ar NVARCHAR(MAX) NULL,

    incident_scope_title_en NVARCHAR(255) NULL,
    incident_scope_title_ar NVARCHAR(255) NULL,

    incident_sensitivity_classification_title_en NVARCHAR(255) NULL,
    incident_sensitivity_classification_title_ar NVARCHAR(255) NULL,

    confidentiality_severity_level_title_en NVARCHAR(255) NULL,
    confidentiality_severity_level_title_ar NVARCHAR(255) NULL,
    confidentiality_severity_level_color_code NVARCHAR(255) NULL,

    integrity_severity_level_title_en NVARCHAR(255) NULL,
    integrity_severity_level_title_ar NVARCHAR(255) NULL,
    integrity_severity_level_color_code NVARCHAR(255) NULL,

    availability_severity_level_title_en NVARCHAR(255) NULL,
    availability_severity_level_title_ar NVARCHAR(255) NULL,
    availability_severity_level_color_code NVARCHAR(255) NULL,

    acknowledged_at DATETIME NULL,
    acknowledged_by_name_en NVARCHAR(255) NULL,
    acknowledged_by_name_ar NVARCHAR(255) NULL,

    closure_resolution_summary NVARCHAR(MAX) NULL,
    closure_notes NVARCHAR(MAX) NULL,
    closure_submitted_at DATETIME NULL,
    closure_submitted_by_name_en NVARCHAR(255) NULL,
    closure_submitted_by_name_ar NVARCHAR(255) NULL,
    closure_approved_at DATETIME NULL,
    closure_approved_by_name_en NVARCHAR(255) NULL,
    closure_approved_by_name_ar NVARCHAR(255) NULL,
    closure_workflow_status_title_en NVARCHAR(255) NULL,
    closure_workflow_status_title_ar NVARCHAR(255) NULL,
    closure_workflow_status_color_code NVARCHAR(255) NULL;

ALTER TABLE incident_report_incidents
ADD
    is_acknowledge BIT NOT NULL
        CONSTRAINT DF_incident_report_incidents_is_acknowledge DEFAULT (0);

ALTER TABLE incident_report_incidents
ADD
    is_sla_time_achieved BIT NOT NULL
        CONSTRAINT DF_incident_report_incidents_is_sla_time_achieved DEFAULT (0);

ALTER TABLE incident_report_incidents
ADD
    is_sla_time_exceeded BIT NOT NULL
        CONSTRAINT DF_incident_report_incidents_is_sla_time_exceeded DEFAULT (0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE incident_report_incidents
DROP CONSTRAINT DF_incident_report_incidents_is_acknowledge;

ALTER TABLE incident_report_incidents
DROP CONSTRAINT DF_incident_report_incidents_is_sla_time_achieved;

ALTER TABLE incident_report_incidents
DROP CONSTRAINT DF_incident_report_incidents_is_sla_time_exceeded;

ALTER TABLE incident_report_incidents
DROP COLUMN
    description,
    quick_action,
    incident_type_title_en,
    incident_type_title_ar,
    incident_status_color_code,
    responsible_unit_title_en,
    responsible_unit_title_ar,
    responsible_unit_color_code,
    business_owner_name_en,
    business_owner_name_ar,
    response_team_lead_name_en,
    response_team_lead_name_ar,
    response_team_members_en,
    response_team_members_ar,
    incident_scope_title_en,
    incident_scope_title_ar,
    incident_sensitivity_classification_title_en,
    incident_sensitivity_classification_title_ar,
    confidentiality_severity_level_title_en,
    confidentiality_severity_level_title_ar,
    confidentiality_severity_level_color_code,
    integrity_severity_level_title_en,
    integrity_severity_level_title_ar,
    integrity_severity_level_color_code,
    availability_severity_level_title_en,
    availability_severity_level_title_ar,
    availability_severity_level_color_code,
    is_acknowledge,
    acknowledged_at,
    acknowledged_by_name_en,
    acknowledged_by_name_ar,
    is_sla_time_achieved,
    is_sla_time_exceeded,
    closure_resolution_summary,
    closure_notes,
    closure_submitted_at,
    closure_submitted_by_name_en,
    closure_submitted_by_name_ar,
    closure_approved_at,
    closure_approved_by_name_en,
    closure_approved_by_name_ar,
    closure_workflow_status_title_en,
    closure_workflow_status_title_ar,
    closure_workflow_status_color_code;
-- +goose StatementEnd
