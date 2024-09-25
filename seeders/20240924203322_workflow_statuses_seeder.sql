-- +goose Up
-- +goose StatementBegin
INSERT INTO workflow_statuses (title_en, title_ar, color_code, type, status)
VALUES
-- Submitted
('Submitted', 'تم التقديم', '#90DD75', 'submit', 'active'),

-- Revert
('Revert', 'إعادة', '#FF9F45', 'revert', 'active'),

-- Reject
('Reject', 'رفض', '#DC3545', 'reject', 'active'),

-- Approve
('Approve', 'موافقة', '#90DD75', 'approve', 'active'),

-- Published
('Published', 'تم النشر', '#007BFF', 'publish', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE workflow_statuses;
-- +goose StatementEnd
