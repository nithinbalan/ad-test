-- +goose Up
-- +goose StatementBegin
INSERT INTO penetration_test_classifications (title_en, title_ar, type, status)
VALUES ('Top-secret', N'سري للغاية', 'top-secret', 'active'), ('Restricted', N'مقيد', 'restricted', 'active'), ('Public', N'عام', 'public', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE penetration_test_classifications;
-- +goose StatementEnd
