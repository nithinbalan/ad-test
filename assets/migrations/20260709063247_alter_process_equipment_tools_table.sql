-- +goose Up
-- +goose StatementBegin
ALTER TABLE process_equipment_tools
ADD type VARCHAR(20) NOT NULL DEFAULT 'other',
    CONSTRAINT CHK_process_equipment_tools_type CHECK (type IN ('other', 'portable_laptop'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE process_equipment_tools
DROP CONSTRAINT CHK_process_equipment_tools_type;

ALTER TABLE process_equipment_tools
DROP COLUMN type;
-- +goose StatementEnd
