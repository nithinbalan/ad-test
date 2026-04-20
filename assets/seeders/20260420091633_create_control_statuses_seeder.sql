-- +goose Up
-- +goose StatementBegin
INSERT INTO control_statuses (title_en, title_ar, color_code, type, status) VALUES
    ('Needs Improvement', 'يحتاج تحسين', '#DC3545', 'needs_improvement', 'active'),
    ('Doesn''t Need Improvement', 'لا يحتاج تحسين', '#28A745', 'doesnt_need_improvement', 'active');

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM control_statuses
WHERE type IN (
    'needs_improvement',
    'doesnt_need_improvement'
);
-- +goose StatementEnd
