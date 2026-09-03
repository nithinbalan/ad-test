-- +goose Up
-- +goose StatementBegin
INSERT INTO planned_advisory_security_classifications (title_en, title_ar, status) VALUES
    ('Top-secret', N'سري للغاية', 'active'),
    ('Secret', N'سري', 'active'),
    ('Restricted', N'مقيد', 'active'),
    ('Public', N'عام', 'active');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM planned_advisory_security_classifications
WHERE title_en IN ('Top-secret', 'Secret', 'Restricted', 'Public');
-- +goose StatementEnd
