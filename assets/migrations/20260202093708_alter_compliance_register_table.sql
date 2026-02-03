-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_register 
ALTER COLUMN compliance_category_id INT NULL;

ALTER TABLE compliance_register 
ALTER COLUMN compliance_status_id INT NULL;

ALTER TABLE compliance_register 
ALTER COLUMN responsible_unit_id INT NULL;

ALTER TABLE compliance_register 
ALTER COLUMN responsible_user_id INT NULL;

ALTER TABLE compliance_register 
ALTER COLUMN lifecycle_status_id INT NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_register 
ALTER COLUMN compliance_category_id INT NOT NULL;

ALTER TABLE compliance_register 
ALTER COLUMN compliance_status_id INT NOT NULL;

ALTER TABLE compliance_register 
ALTER COLUMN responsible_unit_id INT NOT NULL;

ALTER TABLE compliance_register 
ALTER COLUMN responsible_user_id INT NOT NULL;

ALTER TABLE compliance_register 
ALTER COLUMN lifecycle_status_id INT NOT NULL;
-- +goose StatementEnd
