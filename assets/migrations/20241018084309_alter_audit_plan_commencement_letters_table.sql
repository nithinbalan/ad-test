-- +goose Up
-- +goose StatementBegin
-- Add the audit_plan_commencement_letter_status_id column
ALTER TABLE audit_plan_commencement_letters
ADD audit_plan_commencement_letter_status_id INT NULL;

-- Add the sent_at column (nullable because it's only relevant when the letter is sent)
ALTER TABLE audit_plan_commencement_letters
ADD sent_at DATETIME NULL;

-- Add the foreign key constraint
ALTER TABLE audit_plan_commencement_letters
ADD CONSTRAINT FK_audit_plan_commencement_letters_audit_plan_commencement_letter_statuses
FOREIGN KEY (audit_plan_commencement_letter_status_id) REFERENCES audit_plan_commencement_letter_statuses(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the foreign key constraint first
ALTER TABLE audit_plan_commencement_letters
DROP CONSTRAINT FK_audit_plan_commencement_letters_audit_plan_commencement_letter_statuses;

-- Drop the sent_at column
ALTER TABLE audit_plan_commencement_letters
DROP COLUMN sent_at;

-- Drop the audit_plan_commencement_letter_status_id column
ALTER TABLE audit_plan_commencement_letters
DROP COLUMN audit_plan_commencement_letter_status_id;
-- +goose StatementEnd
