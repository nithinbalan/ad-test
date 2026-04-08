-- +goose Up
-- +goose StatementBegin

ALTER TABLE risk_register
ADD 
    risk_code NVARCHAR(500) NULL,
    registration_date DATETIME NULL,
    division_id INT NULL,
    section_id INT NULL,
    expected_cause NVARCHAR(MAX) NULL,
    consequence NVARCHAR(MAX) NULL,
    is_strategic BIT NOT NULL 
        CONSTRAINT DF_risk_register_is_strategic DEFAULT 0,
    is_operational BIT NOT NULL 
        CONSTRAINT DF_risk_register_is_operational DEFAULT 0,
    is_internal BIT NOT NULL 
        CONSTRAINT DF_risk_register_is_internal DEFAULT 0,
    is_external BIT NOT NULL 
        CONSTRAINT DF_risk_register_is_external DEFAULT 0;

ALTER TABLE risk_register
ADD CONSTRAINT FK_risk_register_section_id 
FOREIGN KEY (section_id) REFERENCES organization_units(id);

ALTER TABLE risk_register
ADD CONSTRAINT FK_risk_register_division_id 
FOREIGN KEY (division_id) REFERENCES organization_units(id);

ALTER TABLE risk_register
ALTER COLUMN risk_category_id INT NULL;

ALTER TABLE risk_register
ALTER COLUMN responsible_unit_id INT NULL;

ALTER TABLE risk_register
ALTER COLUMN responsible_user_id INT NULL;

-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin

-- Drop FK first
ALTER TABLE risk_register
DROP CONSTRAINT FK_risk_register_section_id;

ALTER TABLE risk_register
DROP CONSTRAINT FK_risk_register_division_id;

-- Drop DEFAULT constraints
ALTER TABLE risk_register
DROP CONSTRAINT DF_risk_register_is_strategic;

ALTER TABLE risk_register
DROP CONSTRAINT DF_risk_register_is_operational;

ALTER TABLE risk_register
DROP CONSTRAINT DF_risk_register_is_internal;

ALTER TABLE risk_register
DROP CONSTRAINT DF_risk_register_is_external;

-- Then drop columns
ALTER TABLE risk_register
DROP COLUMN 
    risk_code,
    registration_date,
    division_id,
    section_id,
    expected_cause,
    consequence,
    is_strategic,
    is_operational,
    is_internal,
    is_external;

-- +goose StatementEnd