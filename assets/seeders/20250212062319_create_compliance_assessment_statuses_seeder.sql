-- +goose Up
-- +goose StatementBegin
INSERT INTO compliance_assessment_statuses
    (title_en, title_ar, color_code, type, status)
VALUES
    ('Draft', 'مسودة', '#007BFF', 'draft', 'active'),
    ('In Review', 'قيد المراجعة', '#FF9F45', 'in_review', 'active'),
    ('Approved', 'معتمد', '#90DD75', 'approved', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE compliance_assessment_statuses;
-- +goose StatementEnd
