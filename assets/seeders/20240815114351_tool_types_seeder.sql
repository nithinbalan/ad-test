-- +goose Up
-- +goose StatementBegin
INSERT INTO tool_types (title_en, title_ar, status)
VALUES
    ('Hardware', 'أجهزة', 'active'),
    ('Software', 'برمجيات', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE tool_types;
-- +goose StatementEnd
