-- +goose Up
-- +goose StatementBegin

-- 1. Remove all existing data
DELETE FROM control_types;

-- 2. Insert only required records
INSERT INTO control_types (title_en, title_ar, score, status) VALUES
('Preventive',  'الوقائية',   1, 'active'),
('Corrective', 'التصحيحية',   0, 'active'),
('Detective',  'الكشفية',     0, 'active'),
('Deterrent',  'الردعية',     0, 'active');

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- 1. Remove current seeded data
DELETE FROM control_types;

-- 2. Restore original control types (normal state)
INSERT INTO control_types (title_en, title_ar, score, status) VALUES
('Preventive',     'الوقائية',   0, 'active'),
('Detective',      'الكشفية',     0, 'active'),
('Directive',      'التوجيهية',   0, 'active'),
('Corrective',     'التصحيحية',   0, 'active'),
('Compensating',   'التعويضية',   0, 'active'),
('Physical',       'الفيزيائية',  0, 'active'),
('Administrative', 'الإدارية',    0, 'active'),
('Technical',      'التقنية',     0, 'active'),
('Feedback',       'الارتجاعية',  0, 'active'),
('Predictive',     'التنبؤية',    0, 'active');

-- +goose StatementEnd
