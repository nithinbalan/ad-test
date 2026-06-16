-- +goose Up
-- +goose StatementBegin
INSERT INTO awareness_campaign_types
(title_en, title_ar, type, status)
 VALUES
    ('Phishing Campaign', 'حملة التصيد الاحتيالي', 'phishing-campaign', 'active'),
    ('Awareness Campaign', 'حملة التوعية', 'awareness-campaign', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM awareness_campaign_types
-- +goose StatementEnd
