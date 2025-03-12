-- +goose Up
-- +goose StatementBegin
ALTER TABLE risk_register 
ALTER COLUMN responsible_user_id INT NULL;

ALTER TABLE risk_register 
ALTER COLUMN responsible_unit_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE risk_register 
ALTER COLUMN responsible_user_id INT NOT NULL;

ALTER TABLE risk_register 
ALTER COLUMN responsible_unit_id INT NOT NULL;
-- +goose StatementEnd
