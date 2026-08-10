-- +goose Up
-- +goose StatementBegin
ALTER TABLE review_campaign_statuses
DROP CONSTRAINT CHK_review_campaign_statuses_type;

ALTER TABLE review_campaign_statuses
ADD CONSTRAINT CHK_review_campaign_statuses_type
CHECK (
    type IN (
        'planned',
        'in_progress',
        'submitted',
        'approved',
        'closed',
        'rejected'
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM review_campaign_statuses
WHERE type = 'rejected';

ALTER TABLE review_campaign_statuses
DROP CONSTRAINT CHK_review_campaign_statuses_type;

ALTER TABLE review_campaign_statuses
ADD CONSTRAINT CHK_review_campaign_statuses_type
CHECK (
    type IN (
        'planned',
        'in_progress',
        'submitted',
        'approved',
        'closed'
    )
);
-- +goose StatementEnd
