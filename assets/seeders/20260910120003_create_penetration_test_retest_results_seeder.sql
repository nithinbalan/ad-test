-- +goose Up
-- +goose StatementBegin
INSERT INTO penetration_test_retest_results (title_en, title_ar, type, status)
VALUES ('Pass', N'ناجح', 'pass', 'active'), ('Fail', N'فاشل', 'fail', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE penetration_test_retest_results;
-- +goose StatementEnd
