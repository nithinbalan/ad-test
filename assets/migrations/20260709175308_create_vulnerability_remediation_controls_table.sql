-- +goose Up
-- +goose StatementBegin
CREATE TABLE vulnerability_remediation_controls
(
    id INT PRIMARY KEY IDENTITY(1,1),

    vulnerability_id INT NOT NULL,
    control_id INT NOT NULL,

    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_vulnerability_remediation_controls_vulnerability FOREIGN KEY (vulnerability_id)
        REFERENCES vulnerabilities(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_vulnerability_remediation_controls_control FOREIGN KEY (control_id)
        REFERENCES controls(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_vulnerability_remediation_controls_created_by FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE vulnerability_remediation_controls;
-- +goose StatementEnd