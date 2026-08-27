-- +goose Up
-- +goose StatementBegin
INSERT INTO penetration_test_executors (title_en, title_ar, status)
VALUES ('External Vendor', N'مورد خارجي', 'active'), ('CS Protection Team', N'فريق حماية الأمن السيبراني', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE penetration_test_executors;
-- +goose StatementEnd
