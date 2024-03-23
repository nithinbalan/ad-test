-- +goose Up
-- +goose StatementBegin
INSERT INTO organization_levels (level, title_en, title_ar, is_primary) VALUES
(1, 'Division', 'تقسيم', 0),
(2, 'Department', 'قسم', 1),
(3, 'Section', 'تقسيم', 0),
(4, 'Sub Section', 'قسم', 0);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE organization_levels;
-- +goose StatementEnd
