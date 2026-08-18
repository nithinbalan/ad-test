-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_sla_escalation_matrix (
    id INT PRIMARY KEY IDENTITY(1, 1),
    severity_id INT NOT NULL,
    type VARCHAR(20) NOT NULL,
        CONSTRAINT CHK_incident_sla_escalation_matrix_type CHECK (type IN ('escalation_matrix', 'resolution_matrix')),
    notify_within INT DEFAULT NULL,
    notify_within_unit_id INT DEFAULT NULL,
    notify_to_incident_owner BIT NOT NULL DEFAULT 0,
    notify_to_response_team BIT NOT NULL DEFAULT 0,
    notify_role_id INT DEFAULT NULL,
    escalate_after_incident_owner BIT NOT NULL DEFAULT 0,
    escalate_after_response_team BIT NOT NULL DEFAULT 0,
    escalate_role_id INT DEFAULT NULL,
    close_within INT DEFAULT NULL,
    close_within_unit_id INT DEFAULT NULL,
    close_within_incident_owner BIT NOT NULL DEFAULT 0,
    close_within_response_team BIT NOT NULL DEFAULT 0,
    close_within_role_id INT DEFAULT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_incident_sla_escalation_matrix_status CHECK (status IN ('active', 'inactive')),

    CONSTRAINT FK_incident_sla_escalation_matrix_severity_levels
        FOREIGN KEY (severity_id) REFERENCES severity_levels(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_incident_sla_escalation_matrix_notify_within_unit
        FOREIGN KEY (notify_within_unit_id) REFERENCES business_continuity_duration_units(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_incident_sla_escalation_matrix_notify_role
        FOREIGN KEY (notify_role_id) REFERENCES roles(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_incident_sla_escalation_matrix_escalate_role
        FOREIGN KEY (escalate_role_id) REFERENCES roles(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_incident_sla_escalation_matrix_close_within_unit
        FOREIGN KEY (close_within_unit_id) REFERENCES business_continuity_duration_units(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_incident_sla_escalation_matrix_close_within_role
        FOREIGN KEY (close_within_role_id) REFERENCES roles(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_sla_escalation_matrix;
-- +goose StatementEnd
