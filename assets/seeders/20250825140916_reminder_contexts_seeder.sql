-- +goose Up
-- +goose StatementBegin

INSERT INTO reminder_contexts (type, title_en, title_ar) VALUES
('policy-acknowledgement-due', 'Policy Acknowledgement Due', 'إقرار السياسة المستحق');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM reminder_contexts
WHERE type = 'policy-acknowledgement-due';

-- +goose StatementEnd
