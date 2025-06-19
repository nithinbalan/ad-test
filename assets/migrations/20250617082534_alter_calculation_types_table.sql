-- +goose Up
-- +goose StatementBegin
ALTER TABLE calculation_types
DROP CONSTRAINT CHK_calculation_types_type;

-- ALTER TABLE calculation_types
-- ADD CONSTRAINT CHK_calculation_types_type CHECK (type IN ('direct', 'formula'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- ALTER TABLE calculation_types
-- DROP CONSTRAINT CHK_calculation_types_type;

ALTER TABLE calculation_types
ADD CONSTRAINT CHK_calculation_types_type CHECK (type IN ('sum', 'average', 'max', 'min', 'count', 'custom-formula'));
-- +goose StatementEnd
