-- +goose Up
-- +goose StatementBegin
INSERT INTO control_review_statuses (title_en, title_ar, color_code, type, status) VALUES
('Pending', 'قيد الانتظار', '#00FF00', 'pending', 'active'),
('Reviewed', 'تمت المراجعة', '#808080', 'reviewed', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE control_review_statuses;
-- +goose StatementEnd
