-- +goose Up
-- +goose StatementBegin
ALTER TABLE awareness_campaign_types DROP CONSTRAINT CHK_awareness_campaign_types_type;

ALTER TABLE awareness_campaign_types ADD CONSTRAINT CHK_awareness_campaign_types_type
CHECK (type IN (
    'phishing-campaign',
    'awareness-campaign',
    'workshop-event',
    'others'
));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE awareness_campaign_types DROP CONSTRAINT CHK_awareness_campaign_types_type;

ALTER TABLE awareness_campaign_types ADD CONSTRAINT CHK_awareness_campaign_types_type
CHECK (type IN (
    'phishing-campaign',
    'awareness-campaign',
    'others'
));
-- +goose StatementEnd
