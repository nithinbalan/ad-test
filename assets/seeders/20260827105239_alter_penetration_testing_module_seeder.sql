-- +goose Up
-- +goose StatementBegin
UPDATE modules
SET [order] = 17,
    icon = 'icon-penetration-testing.svg'
WHERE identifier = 'penetration_testing';

UPDATE modules
SET title_en = 'Penetration Testing',
    title_ar = N'اختبار الاختراق',
    [order] = 17,
    route = '/settings/penetration-testing'
WHERE identifier = 'penetration_testing_setting';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE modules
SET [order] = 0,
    icon = NULL
WHERE identifier = 'penetration_testing';

UPDATE modules
SET title_en = 'Penetration Testing Masters',
    title_ar = N'البيانات الرئيسية لاختبار الاختراق',
    [order] = 1,
    route = '/penetration-testing'
WHERE identifier = 'penetration_testing_setting';
-- +goose StatementEnd
