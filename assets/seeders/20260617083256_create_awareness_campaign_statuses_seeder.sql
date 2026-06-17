-- +goose Up
-- +goose StatementBegin
INSERT INTO awareness_campaign_statuses
(title_en, title_ar, type, status)
VALUES
    ('Planned', 'مخطط', 'planned', 'active'),
    ('Submitted', 'مقدم', 'submitted', 'active'),
    ('Cancelled', 'ملغى', 'cancelled', 'active'),
    ('Complete', 'مكتمل', 'complete', 'active'),
    ('Not Started', 'لم يبدأ', 'not_started', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM awareness_campaign_statuses
-- +goose StatementEnd
