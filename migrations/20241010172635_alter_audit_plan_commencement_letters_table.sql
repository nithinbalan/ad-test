-- +goose Up
-- +goose StatementBegin
ALTER TABLE audit_plan_commencement_letters
ADD subject_en NVARCHAR(255) NOT NULL,
    subject_ar NVARCHAR(255) NOT NULL,
    body_en TEXT NOT NULL,
    body_ar TEXT NOT NULL;

ALTER TABLE audit_plan_commencement_letters
DROP COLUMN subject;

ALTER TABLE audit_plan_commencement_letters
DROP COLUMN body;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE audit_plan_commencement_letters
ADD subject NVARCHAR(255) NOT NULL;

ALTER TABLE audit_plan_commencement_letters
ADD body TEXT NOT NULL;

ALTER TABLE audit_plan_commencement_letters
DROP COLUMN subject_en;

ALTER TABLE audit_plan_commencement_letters
DROP COLUMN subject_ar;

ALTER TABLE audit_plan_commencement_letters
DROP COLUMN body_en;

ALTER TABLE audit_plan_commencement_letters
DROP COLUMN body_ar;
-- +goose StatementEnd
