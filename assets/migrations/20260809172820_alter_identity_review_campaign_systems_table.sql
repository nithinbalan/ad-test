-- +goose Up
-- +goose StatementBegin
ALTER TABLE identity_review_campaign_systems
ALTER COLUMN system_register_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM identity_review_campaign_systems
WHERE system_register_id IS NULL;

ALTER TABLE identity_review_campaign_systems
ALTER COLUMN system_register_id INT NOT NULL;
-- +goose StatementEnd
