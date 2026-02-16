-- +goose Up
-- +goose StatementBegin

-- Insert master data
INSERT INTO control_automations (title_en, title_ar, type, score, status, created_at, created_by_id)
VALUES
('Not Automated',        'غير آلي',            'not_automated',        0.00, 'active', GETDATE(), 1),
('Partially Automated',  'آلي جزئياً',         'partially_automated',  0.50, 'active', GETDATE(), 1),
('Fully Automated',      'آلي بالكامل',        'fully_automated',      1.00, 'active', GETDATE(), 1);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Rollback to normal state (empty master table)
DELETE FROM control_automations;

-- +goose StatementEnd
