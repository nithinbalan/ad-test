-- +goose Up
-- +goose StatementBegin
INSERT INTO internal_committee_statuses (title_en, title_ar, color_code, type, status) VALUES
    ('Draft', 'مسودة', '#17A2B8', 'draft', 'active'),
    ('Submitted', 'مُقدَّم', '#007BFF ', 'submitted', 'active'),
    ('In Review', 'قيد المراجعة', '#FFA500', 'in_review', 'active'),
    ('In Progress', 'قيد التنفيذ', '#FD7E14', 'in_progress', 'active'),
    ('Approved', 'معتمد', '#28A745', 'approved', 'active'),
    ('Reverted', 'ارجاع', '#FFC107', 'revert', 'active'),
    ('Reformation', 'مراجعة', '#6C757D', 'reformation', 'active'),
    ('Disbanded', 'ملغى', '#DC3545', 'disbanded', 'active');

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM internal_committee_statuses
WHERE type IN (
    'draft',
    'submitted',
    'in_review',
    'in_progress',
    'approved',
    'revert',
    'reformation',
    'disbanded'
);
-- +goose StatementEnd
