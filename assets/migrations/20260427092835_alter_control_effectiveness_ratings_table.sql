-- +goose Up
-- +goose StatementBegin
ALTER TABLE control_effectiveness_ratings
ADD color_code VARCHAR(25) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE control_effectiveness_ratings
DROP COLUMN color_code;
-- +goose StatementEnd
