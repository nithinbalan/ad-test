-- +goose Up
-- +goose StatementBegin

ALTER TABLE incident_register
ADD business_owner_id INT NULL;

ALTER TABLE incident_register
ADD incident_scope_id INT NULL;

ALTER TABLE incident_register
ADD incident_sensitivity_classification_id INT NULL;

ALTER TABLE incident_register
ADD confidentiality_severity_level_id INT NULL;

ALTER TABLE incident_register
ADD integrity_severity_level_id INT NULL;

ALTER TABLE incident_register
ADD availability_severity_level_id INT NULL;

ALTER TABLE incident_register
ADD incident_response_team_id INT NULL;

ALTER TABLE incident_register
ADD response_team_lead_id INT NULL;

ALTER TABLE incident_register
ADD
    is_acknowledge BIT NOT NULL
        CONSTRAINT DF_incident_register_is_acknowledge DEFAULT (0);

ALTER TABLE incident_register
ADD acknowledged_at DATETIME NULL;

ALTER TABLE incident_register
ADD acknowledged_by_id INT NULL;

ALTER TABLE incident_register
ADD scheduler_last_run_at DATETIME NULL;

ALTER TABLE incident_register
ADD
    is_sla_time_achieved BIT NOT NULL
        CONSTRAINT DF_incident_register_is_sla_time_achieved DEFAULT (0);

ALTER TABLE incident_register
ADD
    is_sla_time_exceeded BIT NOT NULL
        CONSTRAINT DF_incident_register_is_sla_time_exceeded DEFAULT (0);

ALTER TABLE incident_register
ADD CONSTRAINT FK_incident_register_business_owner
    FOREIGN KEY (business_owner_id) REFERENCES users (id)
    ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE incident_register
ADD CONSTRAINT FK_incident_register_incident_scopes
    FOREIGN KEY (incident_scope_id) REFERENCES incident_scopes (id)
    ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE incident_register
ADD CONSTRAINT FK_incident_register_incident_sensitivity_classifications
    FOREIGN KEY (incident_sensitivity_classification_id) REFERENCES incident_sensitivity_classifications (id)
    ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE incident_register
ADD CONSTRAINT FK_incident_register_confidentiality_severity_levels
    FOREIGN KEY (confidentiality_severity_level_id) REFERENCES severity_levels (id)
    ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE incident_register
ADD CONSTRAINT FK_incident_register_integrity_severity_levels
    FOREIGN KEY (integrity_severity_level_id) REFERENCES severity_levels (id)
    ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE incident_register
ADD CONSTRAINT FK_incident_register_availability_severity_levels
    FOREIGN KEY (availability_severity_level_id) REFERENCES severity_levels (id)
    ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE incident_register
ADD CONSTRAINT FK_incident_register_incident_response_team
    FOREIGN KEY (incident_response_team_id) REFERENCES incident_response_teams (id)
    ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE incident_register
ADD CONSTRAINT FK_incident_register_response_team_lead
    FOREIGN KEY (response_team_lead_id) REFERENCES users (id)
    ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE incident_register
ADD CONSTRAINT FK_incident_register_acknowledged_by
    FOREIGN KEY (acknowledged_by_id) REFERENCES users (id)
    ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE incident_register
DROP CONSTRAINT DF_incident_register_is_acknowledge;

ALTER TABLE incident_register
DROP CONSTRAINT DF_incident_register_is_sla_time_achieved;

ALTER TABLE incident_register
DROP CONSTRAINT DF_incident_register_is_sla_time_exceeded;

ALTER TABLE incident_register
DROP CONSTRAINT FK_incident_register_availability_severity_levels;

ALTER TABLE incident_register
DROP CONSTRAINT FK_incident_register_integrity_severity_levels;

ALTER TABLE incident_register
DROP CONSTRAINT FK_incident_register_confidentiality_severity_levels;

ALTER TABLE incident_register
DROP CONSTRAINT FK_incident_register_incident_sensitivity_classifications;

ALTER TABLE incident_register
DROP CONSTRAINT FK_incident_register_incident_scopes;

ALTER TABLE incident_register
DROP CONSTRAINT FK_incident_register_business_owner;

ALTER TABLE incident_register
DROP CONSTRAINT FK_incident_register_incident_response_team;

ALTER TABLE incident_register
DROP CONSTRAINT FK_incident_register_response_team_lead;

ALTER TABLE incident_register
DROP CONSTRAINT FK_incident_register_acknowledged_by;

ALTER TABLE incident_register
DROP COLUMN availability_severity_level_id;

ALTER TABLE incident_register
DROP COLUMN integrity_severity_level_id;

ALTER TABLE incident_register
DROP COLUMN confidentiality_severity_level_id;

ALTER TABLE incident_register
DROP COLUMN incident_sensitivity_classification_id;

ALTER TABLE incident_register
DROP COLUMN incident_response_team_id;

ALTER TABLE incident_register
DROP COLUMN response_team_lead_id;

ALTER TABLE incident_register
DROP COLUMN is_acknowledge;

ALTER TABLE incident_register
DROP COLUMN acknowledged_at;

ALTER TABLE incident_register
DROP COLUMN scheduler_last_run_at;

ALTER TABLE incident_register
DROP COLUMN acknowledged_by_id;

ALTER TABLE incident_register
DROP COLUMN incident_scope_id;

ALTER TABLE incident_register
DROP COLUMN business_owner_id;

-- +goose StatementEnd
