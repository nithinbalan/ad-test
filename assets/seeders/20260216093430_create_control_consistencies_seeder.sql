-- +goose Up
-- +goose StatementBegin

-- Insert master data
INSERT INTO control_consistencies (title_en, title_ar, type, score, status, created_at, created_by_id)
VALUES
('Consistent',   'متسق',        'consistent',   1.00, 'active', GETDATE(), 1),
('Inconsistent', 'غير متسق',    'inconsistent', 0.50, 'active', GETDATE(), 1),
('Not Applied',  'غير مطبق',    'not_applied',  0.00, 'active', GETDATE(), 1);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Rollback to normal state (empty master table)
DELETE FROM control_consistencies;

-- +goose StatementEnd
