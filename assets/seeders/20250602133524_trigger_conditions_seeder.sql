-- +goose Up
-- +goose StatementBegin
INSERT INTO trigger_conditions (title_en, title_ar,type, status) VALUES
('Before Due Date', 'قبل تاريخ الاستحقاق','before-due-date', 'active'),
('On Due Date', 'في تاريخ الاستحقاق', 'on-due-date', 'active'),
('After Due Date', 'بعد تاريخ الاستحقاق', 'after-due-date', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE trigger_conditions;
-- +goose StatementEnd