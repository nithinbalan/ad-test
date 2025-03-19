-- +goose Up
-- +goose StatementBegin
ALTER TABLE time_units 
ADD type VARCHAR(25) NULL,
    CONSTRAINT CHK_time_units_type CHECK (type IN ('hour', 'minute', 'day', 'week'));

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE time_units 
DROP COLUMN type;
-- +goose StatementEnd
