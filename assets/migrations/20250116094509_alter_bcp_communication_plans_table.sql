-- +goose Up
-- +goose StatementBegin
-- Add the bcp_process_id column with a foreign key to bcp_communication_plans table

ALTER TABLE bcp_communication_plans
ADD bcp_process_id INT NULL;

ALTER TABLE bcp_communication_plans
ADD CONSTRAINT FK_bcp_communication_plans_bcp_processes
FOREIGN KEY (bcp_process_id) REFERENCES bcp_processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the foreign key constraint for bcp_process_id
ALTER TABLE bcp_communication_plans
DROP CONSTRAINT FK_bcp_communication_plans_bcp_processes;

-- Drop the bcp_process_id columns
ALTER TABLE bcp_communication_plans
DROP COLUMN bcp_process_id;
-- +goose StatementEnd