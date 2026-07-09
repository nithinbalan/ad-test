-- +goose Up
-- +goose StatementBegin
CREATE TABLE sla_escalation_matrix (
    id INT PRIMARY KEY IDENTITY(1,1),
    severity_id INT NOT NULL,
    remediation_days INT NOT NULL,
    is_email BIT NOT NULL DEFAULT 0,
    is_alert BIT NOT NULL DEFAULT 0,
    role_id INT NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',

    CONSTRAINT CHK_sla_escalation_matrix_status
        CHECK (status IN ('active', 'inactive')),
    CONSTRAINT FK_sla_escalation_matrix_severity_level FOREIGN KEY (severity_id)
      REFERENCES severity_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_escalation_matrix_role FOREIGN KEY (role_id)
      REFERENCES roles(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sla_escalation_matrix;
-- +goose StatementEnd

