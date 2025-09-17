-- +goose Up
-- +goose StatementBegin
INSERT INTO comment_statuses (title_en, title_ar, color_code, type, status) VALUES
('Open', 'يفتح', '#007BFF', 'open', 'active'),
('Resolved', 'تم الحل', '#28A745', 'resolved', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE comment_statuses;
-- +goose StatementEnd
