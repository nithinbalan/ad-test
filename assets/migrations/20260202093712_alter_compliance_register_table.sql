-- +goose Up
-- +goose StatementBegin

ALTER TABLE compliance_register
ADD 
    legislative_body_id INT NULL,
    legislation_type_id INT NULL,
    legislation_no NVARCHAR(500) NULL,
    importance_classification_id INT NULL,
    secrecy_classification_id INT NULL,
    subject NVARCHAR(MAX) NULL,
    is_applicable BIT NOT NULL CONSTRAINT df_compliance_register_is_applicable DEFAULT 0,
    is_not_applicable BIT NOT NULL CONSTRAINT df_compliance_register_is_not_applicable DEFAULT 0,
    is_regulation_notification_assessment BIT NOT NULL CONSTRAINT df_compliance_register_is_regulation_notification_assessment DEFAULT 0,
    is_risk_based_assessment BIT NOT NULL CONSTRAINT df_compliance_register_is_risk_based_assessment DEFAULT 0,
    is_public BIT NOT NULL CONSTRAINT df_compliance_register_is_public DEFAULT 0,
    is_shared BIT NOT NULL CONSTRAINT df_compliance_register_is_shared DEFAULT 0;

ALTER TABLE compliance_register
ADD CONSTRAINT FK_compliance_register_legislative_body 
FOREIGN KEY (legislative_body_id) 
REFERENCES compliance_legislative_bodies(id);

ALTER TABLE compliance_register
ADD CONSTRAINT FK_compliance_register_legislation_type 
FOREIGN KEY (legislation_type_id) 
REFERENCES compliance_legislation_types(id);

ALTER TABLE compliance_register
ADD CONSTRAINT FK_compliance_register_importance_classification 
FOREIGN KEY (importance_classification_id) 
REFERENCES compliance_importance_classifications(id);

ALTER TABLE compliance_register
ADD CONSTRAINT FK_compliance_register_secrecy_classification 
FOREIGN KEY (secrecy_classification_id) 
REFERENCES compliance_secrecy_classifications(id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE compliance_register   
DROP CONSTRAINT FK_compliance_register_legislative_body;

ALTER TABLE compliance_register
DROP CONSTRAINT FK_compliance_register_legislation_type;

ALTER TABLE compliance_register
DROP CONSTRAINT FK_compliance_register_importance_classification;

ALTER TABLE compliance_register
DROP CONSTRAINT FK_compliance_register_secrecy_classification;

ALTER TABLE compliance_register
DROP CONSTRAINT df_compliance_register_is_applicable;

ALTER TABLE compliance_register
DROP CONSTRAINT df_compliance_register_is_not_applicable;

ALTER TABLE compliance_register
DROP CONSTRAINT df_compliance_register_is_regulation_notification_assessment;

ALTER TABLE compliance_register
DROP CONSTRAINT df_compliance_register_is_risk_based_assessment;

ALTER TABLE compliance_register
DROP CONSTRAINT df_compliance_register_is_public;

ALTER TABLE compliance_register
DROP CONSTRAINT df_compliance_register_is_shared;

ALTER TABLE compliance_register
DROP COLUMN 
    legislative_body_id,
    legislation_type_id,
    legislation_no,
    importance_classification_id,
    secrecy_classification_id,
    subject,
    is_applicable,
    is_not_applicable,
    is_regulation_notification_assessment,
    is_risk_based_assessment,
    is_public,
    is_shared;

-- +goose StatementEnd
