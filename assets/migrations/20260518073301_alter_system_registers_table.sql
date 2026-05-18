-- +goose Up
-- +goose StatementBegin

ALTER TABLE system_registers
ADD 
    title_en NVARCHAR(1000) NULL,
    title_ar NVARCHAR(1000) NULL;



-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin
-- Then drop columns
ALTER TABLE system_registers
DROP COLUMN 
    title_en,
    title_ar;

-- +goose StatementEnd
