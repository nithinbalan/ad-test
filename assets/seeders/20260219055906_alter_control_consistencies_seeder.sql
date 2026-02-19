-- +goose Up
-- +goose StatementBegin

UPDATE control_consistencies
SET 
    description_en = 'The control is implemented according to a clear schedule (daily, weekly, monthly, etc.) and its implementation is continuously documented.',
    description_ar = 'يتم تنفيذ الضابط وفق جدول زمني واضح (يومي، أسبوعي، شهري، إلخ) ويتم توثيق تنفيذه بشكل مستمر.'
WHERE type = 'consistent';

UPDATE control_consistencies
SET 
    description_en = 'It is applied only when cases occur, or without a clear timetable, or at intermittent intervals.',
    description_ar = 'يتم تطبيقه فقط عند حدوث الحالات، أو بدون جدول زمني واضح، أو على فترات متقطعة.'
WHERE type = 'inconsistent';

UPDATE control_consistencies
SET 
    description_en = 'The control is purely theoretical, or there is no evidence of its actual application.',
    description_ar = 'الضابط نظري فقط، أو لا يوجد دليل على تطبيقه الفعلي.'
WHERE type = 'not_applied';

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

UPDATE control_consistencies
SET 
    description_en = NULL,
    description_ar = NULL
WHERE type IN ('consistent', 'inconsistent', 'not_applied');

-- +goose StatementEnd
