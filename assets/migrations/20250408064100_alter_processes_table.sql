-- +goose Up
-- +goose StatementBegin
ALTER TABLE processes
ADD 
    impact_justification NVARCHAR(MAX) NULL,
    process_criticality_level_id INT NULL,
    criticality_justification NVARCHAR(MAX) NULL,
    CONSTRAINT FK_processes_process_criticality_levels FOREIGN KEY (process_criticality_level_id) REFERENCES process_criticality_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE processes
DROP CONSTRAINT FK_processes_process_criticality_levels,
DROP COLUMN impact_justification,
DROP COLUMN criticality_justification,
DROP COLUMN process_criticality_id;
-- +goose StatementEnd
