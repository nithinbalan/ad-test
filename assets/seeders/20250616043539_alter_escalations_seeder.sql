-- +goose Up
-- +goose StatementBegin
UPDATE escalations
SET description_en = 'Escalate if the expiry date is reached without action.'
WHERE title_en = 'Compliance Register Expiry';

UPDATE escalations
SET description_ar = 'تصعيد إذا تم الوصول إلى تاريخ انتهاء الصلاحية دون اتخاذ أي إجراء.'
WHERE title_en = 'Compliance Register Expiry';

UPDATE escalations
SET description_en = 'Escalate if the compliance deadline is today or missed.'
WHERE title_en = 'Compliance Requirement Deadline';

UPDATE escalations
SET description_ar = 'قم بالتصعيد إذا كان الموعد النهائي للامتثال هو اليوم أو تم تفويته.'
WHERE title_en = 'Compliance Requirement Deadline';

UPDATE escalations
SET description_en = 'Escalate if the action plan is not completed by the target date.'
WHERE title_en = 'Compliance Action Plan Overdue';

UPDATE escalations
SET description_ar = 'تصعيد الأمر إذا لم يتم إكمال خطة العمل بحلول الموعد المستهدف.'
WHERE title_en = 'Compliance Action Plan Overdue';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
