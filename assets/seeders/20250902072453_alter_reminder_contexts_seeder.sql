-- +goose Up
-- +goose StatementBegin

INSERT INTO reminder_contexts (type, title_en, title_ar) VALUES
('compliance-review-due', 'Compliance Review Date Due', 'تاريخ استحقاق مراجعة الامتثال');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM reminder_contexts
WHERE type = 'compliance-review-date-due';

-- +goose StatementEnd
