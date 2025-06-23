-- +goose Up
-- +goose StatementBegin
INSERT INTO audit_plan_commencement_letter_statuses
    (title_en, title_ar, color_code, type, status)
VALUES
    ('Draft', 'مسودة', '#FFC107', 'draft', 'active'),
    ('Sent', 'مرسل', '#28A745', 'sent', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_plan_commencement_letter_statuses;
-- +goose StatementEnd
