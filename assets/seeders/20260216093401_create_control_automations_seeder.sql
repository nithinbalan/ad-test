-- +goose Up
-- +goose StatementBegin

-- Insert master data
INSERT INTO control_automations (title_en, title_ar, type, score, status)
VALUES
('Not Automated',        'غير آلي',            'not_automated',        0.00, 'active'),
('Partially Automated',  'آلي جزئياً',         'partially_automated',  0.50, 'active'),
('Fully Automated',      'آلي بالكامل',        'fully_automated',      1.00, 'active');

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Rollback to normal state (empty master table)
DELETE FROM control_automations;

-- +goose StatementEnd
