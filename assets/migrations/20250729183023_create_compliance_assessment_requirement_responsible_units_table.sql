-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_assessment_requirement_responsible_units (
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_assessment_requirement_id INT NOT NULL,
    responsible_unit_id INT NOT NULL,
    deadline DATE NULL,
    compliance_assessment_framework_value_id INT NULL,
    score INT NULL,
    comments NVARCHAR(MAX) NULL,
    is_completed TINYINT NOT NULL DEFAULT 0,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_carru_requirement FOREIGN KEY (compliance_assessment_requirement_id)
        REFERENCES compliance_assessment_requirements(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_carru_unit FOREIGN KEY (responsible_unit_id)
        REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_carru_framework FOREIGN KEY (compliance_assessment_framework_value_id)
        REFERENCES compliance_assessment_framework_values(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_carru_created_by FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_carru_updated_by FOREIGN KEY (updated_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_assessment_requirement_responsible_units;
-- +goose StatementEnd
