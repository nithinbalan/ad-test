-- +goose Up
-- +goose StatementBegin
ALTER TABLE asset_compliance_requirements
ADD compliance_register_id INT NOT NULL;

ALTER TABLE asset_compliance_requirements
ADD CONSTRAINT FK_asset_compliance_requirement_compliance_register
FOREIGN KEY (compliance_register_id)
REFERENCES compliance_register(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE asset_compliance_requirements
DROP CONSTRAINT FK_asset_compliance_requirement_compliance_register;

-- Drop the compliance_register_id column
ALTER TABLE asset_compliance_requirements
DROP COLUMN compliance_register_id;

-- +goose StatementEnd
