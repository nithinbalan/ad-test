-- +goose Up
-- +goose StatementBegin

-- Rename single recovery objective columns to *_from_value
EXEC sp_rename 'bcp_processes.rto', 'rto_from_value', 'COLUMN';
EXEC sp_rename 'bcp_processes.mtpd', 'mtpd_from_value', 'COLUMN';
EXEC sp_rename 'bcp_processes.rpo', 'rpo_from_value', 'COLUMN';

-- Add *_to_value columns
ALTER TABLE bcp_processes
ADD rto_to_value INT NULL,
    mtpd_to_value INT NULL,
    rpo_to_value INT NULL;

-- Add RSP recovery start point columns
ALTER TABLE bcp_processes
ADD rsp_from_value INT NULL,
    rsp_to_value INT NULL,
    rsp_condition VARCHAR(20) NULL,
    rsp_unit_id INT NULL;

-- Add FK for rsp_unit_id
ALTER TABLE bcp_processes
ADD CONSTRAINT FK_bcp_processes_rsp_bc_duration_units
FOREIGN KEY (rsp_unit_id)
REFERENCES business_continuity_duration_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop RSP FK and columns
ALTER TABLE bcp_processes
DROP CONSTRAINT FK_bcp_processes_rsp_bc_duration_units;

ALTER TABLE bcp_processes
DROP COLUMN rsp_from_value, rsp_to_value, rsp_condition, rsp_unit_id;

-- Drop *_to_value columns
ALTER TABLE bcp_processes
DROP COLUMN rto_to_value, mtpd_to_value, rpo_to_value;

-- Rename columns back
EXEC sp_rename 'bcp_processes.rto_from_value', 'rto', 'COLUMN';
EXEC sp_rename 'bcp_processes.mtpd_from_value', 'mtpd', 'COLUMN';
EXEC sp_rename 'bcp_processes.rpo_from_value', 'rpo', 'COLUMN';

-- +goose StatementEnd
