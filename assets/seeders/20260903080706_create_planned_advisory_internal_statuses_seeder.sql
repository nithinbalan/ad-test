-- +goose Up
-- +goose StatementBegin
INSERT INTO planned_advisory_internal_statuses (title_en, title_ar, status) VALUES
    ('Secondary', N'ثانوي', 'active'),
    ('Internal', N'داخلي', 'active'),
    ('Processing', N'قيد المعالجة', 'active');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM planned_advisory_internal_statuses
WHERE title_en IN ('Secondary', 'Internal', 'Processing');
-- +goose StatementEnd
