-- +goose Up
-- +goose StatementBegin
INSERT INTO importance_classifications (
    type,
    title_en,
    title_ar,
    description_en,
    description_ar,
    status,
    created_at
) VALUES
(
    'normal',
    'Normal',
    N'عادي',
    'Normal importance level',
    N'مستوى أهمية عادي',
    'active',
    GETDATE()
),
(
    'urgent',
    'Urgent',
    N'عاجل',
    'Urgent importance level',
    N'مستوى أهمية عاجل',
    'active',
    GETDATE()
),
(
    'very_urgent',
    'Very Urgent',
    N'عاجل جداً',
    'Very urgent importance level',
    N'مستوى أهمية عاجل جداً',
    'active',
    GETDATE()
),
(
    'important_and_very_urgent',
    'Important and Very Urgent',
    N'مهم وعاجل جداً',
    'Important and very urgent importance level',
    N'مستوى أهمية مهم وعاجل جداً',
    'active',
    GETDATE()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM importance_classifications
WHERE type IN (
    'normal',
    'urgent',
    'very_urgent',
    'important_and_very_urgent'
);
-- +goose StatementEnd
