-- +goose Up
-- +goose StatementBegin
ALTER TABLE control_types
ADD score INT NOT NULL
    CONSTRAINT DF_control_types_score DEFAULT 0;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE control_types
DROP CONSTRAINT DF_control_types_score;

ALTER TABLE control_types
DROP COLUMN score;
-- +goose StatementEnd
