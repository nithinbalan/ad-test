-- +goose Up
-- +goose StatementBegin
INSERT INTO compliance_secrecy_classifications (
    type,
    title_en,
    title_ar,
    description_en,
    description_ar,
    status
) VALUES
(
    'normal',
    'Normal',
    N'عادي',
    'Normal secrecy level',
    N'مستوى سرية عادي',
    'active'
),
(
    'secret',
    'Secret',
    N'سري',
    'Secret secrecy level',
    N'مستوى سرية سري',
    'active'
),
(
    'top_secret',
    'Top Secret',
    N'سري للغاية',
    'Top secret secrecy level',
    N'مستوى سرية سري للغاية',
    'active'
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DELETE FROM compliance_secrecy_classifications
WHERE type IN (
    'normal',
    'secret',
    'top_secret'
);
-- +goose StatementEnd
