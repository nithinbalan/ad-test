-- +goose Up
-- +goose StatementBegin
INSERT INTO organization_levels (level, title_en, title_ar, is_primary, type) VALUES
(1, 'Division', 'تقسيم', 0, 'division'),
(2, 'Department', 'قسم', 1, 'department'),
(3, 'Section', 'تقسيم', 0, 'section'),
(4, 'Sub Section', 'قسم', 0, 'subsection');


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE organization_levels;
-- +goose StatementEnd
