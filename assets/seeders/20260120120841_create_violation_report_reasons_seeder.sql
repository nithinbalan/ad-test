-- +goose Up
-- +goose StatementBegin
INSERT INTO violation_report_reasons
    (title_en, title_ar, description_en, description_ar, status)
VALUES
    ('Risk', 'مخاطر', NULL, NULL, 'active'),
    ('Unfairness', 'عدم العدالة', NULL, NULL, 'active');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM violation_report_reasons
WHERE title_en IN ('Risk', 'Unfairness');
-- +goose StatementEnd
