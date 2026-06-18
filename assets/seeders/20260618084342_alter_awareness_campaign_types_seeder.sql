-- +goose Up
-- +goose StatementBegin
INSERT INTO awareness_campaign_types (title_en, title_ar, type, status)
VALUES
    ('Workshop-Event', N'ورشة عمل / فعالية', 'workshop-event', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM awareness_campaign_types WHERE type = 'workshop-event';
-- +goose StatementEnd
