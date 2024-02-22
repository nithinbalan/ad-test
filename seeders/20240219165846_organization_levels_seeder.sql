-- +goose Up
-- +goose StatementBegin
INSERT INTO organization_levels (level, title_en, title_ar, is_primary, type) VALUES
(1, 'Division', 'تقسيم', 0, 'division'),
(2, 'Department', 'قسم', 1, 'department');


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE organization_levels;
-- +goose StatementEnd
