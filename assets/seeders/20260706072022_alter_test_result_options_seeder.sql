-- +goose Up
-- +goose StatementBegin
INSERT INTO test_result_options (title_en, title_ar, color_code, type, status)
VALUES
    ('Not Tested', N'لم يتم الاختبار', '#FD7E14', 'not-tested', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM test_result_options
WHERE type IN (
    'not-tested'
);
-- +goose StatementEnd
