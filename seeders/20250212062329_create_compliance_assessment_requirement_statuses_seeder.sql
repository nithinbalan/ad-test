-- +goose Up
-- +goose StatementBegin
INSERT INTO compliance_assessment_requirement_statuses
    (title_en, title_ar, color_code, type, status)
VALUES
    ('Approved', 'معتمد', '#90DD75', 'approved', 'active'),
    ('Reverted', 'تم الإرجاع', '#FF9F45', 'reverted', 'active'),
    ('Resubmitted', 'أعيد الإرسال', '#007BFF', 'resubmitted', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE compliance_assessment_requirement_statuses;
-- +goose StatementEnd
