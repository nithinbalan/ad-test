-- +goose Up
-- +goose StatementBegin
ALTER TABLE processes 
ADD risk_level_id INT NULL;

-- Add the foreign key constraints for risk_level_id
ALTER TABLE processes
ADD CONSTRAINT FK_processes_risk_levels FOREIGN KEY (risk_level_id) REFERENCES risk_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop foreign key constraints before removing columns
ALTER TABLE processes
DROP CONSTRAINT FK_processes_risk_levels;

ALTER TABLE processes
DROP COLUMN risk_level_id;

-- +goose StatementEnd
