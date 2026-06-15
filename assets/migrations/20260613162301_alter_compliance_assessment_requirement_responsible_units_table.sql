-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_assessment_requirement_responsible_units
ADD is_acknowledged BIT NOT NULL
    CONSTRAINT DF_compliance_assessment_requirement_responsible_units_is_acknowledged DEFAULT 0;

ALTER TABLE compliance_assessment_requirement_responsible_units
ADD acknowledge_user_id INT NULL;

ALTER TABLE compliance_assessment_requirement_responsible_units
ADD acknowledge_comment NVARCHAR(MAX) NULL;

ALTER TABLE compliance_assessment_requirement_responsible_units
ADD CONSTRAINT FK_carru_acknowledge_user FOREIGN KEY (acknowledge_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_assessment_requirement_responsible_units
DROP CONSTRAINT IF EXISTS FK_carru_acknowledge_user;

ALTER TABLE compliance_assessment_requirement_responsible_units
DROP CONSTRAINT IF EXISTS DF_compliance_assessment_requirement_responsible_units_is_acknowledged;

ALTER TABLE compliance_assessment_requirement_responsible_units
DROP COLUMN IF EXISTS acknowledge_comment;

ALTER TABLE compliance_assessment_requirement_responsible_units
DROP COLUMN IF EXISTS acknowledge_user_id;

ALTER TABLE compliance_assessment_requirement_responsible_units
DROP COLUMN IF EXISTS is_acknowledged;
-- +goose StatementEnd
