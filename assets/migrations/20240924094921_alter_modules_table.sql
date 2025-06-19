-- +goose Up
-- +goose StatementBegin
ALTER TABLE modules 
ADD icon VARCHAR(45),
    route VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE modules 
DROP COLUMN IF EXISTS icon, 
    DROP COLUMN IF EXISTS route;
-- +goose StatementEnd
