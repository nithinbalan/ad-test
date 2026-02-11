-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_assessment_residual_controls 
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_assessment_requirement_responsible_unit_id INT NOT NULL,
    compliance_assessment_impact_type_id INT NOT NULL,
    control_id INT NULL,
    link_control BIT NOT NULL DEFAULT 0,
    is_not_applicable BIT NOT NULL DEFAULT 0,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_compliance_assessment_residual_controls_compliance_assessment_requirement_responsible_unit FOREIGN KEY (compliance_assessment_requirement_responsible_unit_id) REFERENCES compliance_assessment_requirement_responsible_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_residual_controls_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_residual_controls_impact_types FOREIGN KEY (compliance_assessment_impact_type_id) REFERENCES compliance_assessment_impact_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_residual_controls_controls FOREIGN KEY (control_id) REFERENCES controls(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_assessment_residual_controls;
-- +goose StatementEnd
