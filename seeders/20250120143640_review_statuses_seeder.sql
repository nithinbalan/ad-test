-- +goose Up
-- +goose StatementBegin
INSERT INTO review_statuses (title_en, title_ar, color_code, type, status) VALUES
    ('Review Pending', 'قيد المراجعة', '#90DD75', 'review-pending', 'active'),
    ('Reviewed', 'تمت المراجعة', '#6C757D', 'reviewed', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM review_statuses
WHERE type IN ('review-pending', 'reviewed');
-- +goose StatementEnd
