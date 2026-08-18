-- +goose Up
-- +goose StatementBegin

UPDATE incident_statuses
SET
    title_en = 'New',
    title_ar = N'جديد',
    description_en = 'The incident has been newly created.',
    description_ar = N'تم إنشاء الحادث حديثاُ.',
    color_code = '#17A2B8',
    type = 'new'
WHERE type = 'open';

INSERT INTO incident_statuses
    (title_en, title_ar, description_en, description_ar, color_code, type, status)
SELECT * FROM (VALUES
    ('Triage', N'فرز', 'The incident is being triaged and prioritized.', N'يتم فرز الحادث وتحديد أولوياته.', '#20C997', 'triage', 'active'),
    ('Investigating', N'قيد التحقيق', 'The incident is under investigation.', N'الحادث قيد التحقيق.', '#007BFF', 'investigating', 'active'),
    ('Containment', N'احتواء', 'The incident is being contained.', N'يتم احتواء الحادث.', '#FD7E14', 'containment', 'active'),
    ('Eradication', N'استئصال', 'The incident is being eradicated.', N'يتم استئصال الحادث.', '#DC3545', 'eradication', 'active'),
    ('Recovery', N'استرداد', 'The incident is in recovery phase.', N'الحادث في مرحلة الاسترداد.', '#FFC107', 'recovery', 'active'),
    ('Closure Preparation', N'إعداد الإغلاق', 'The incident is in closure preparation phase.', N'الحادث في مرحلة إعداد الإغلاق.', '#FFA500', 'closure_preparation', 'active')
) AS v(title_en, title_ar, description_en, description_ar, color_code, type, status)
WHERE NOT EXISTS (
    SELECT 1 FROM incident_statuses s WHERE s.type = v.type
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

UPDATE incident_statuses
SET
    title_en = 'Open',
    title_ar = N'مفتوح',
    description_en = 'The incident is open and under investigation.',
    description_ar = N'الحادث مفتوح وقيد التحقيق.',
    color_code = '#FD7E14',
    type = 'open'
WHERE type = 'new';

DELETE FROM incident_statuses
WHERE type IN ('triage', 'investigating', 'containment', 'eradication', 'recovery', 'closure_preparation');

-- +goose StatementEnd
