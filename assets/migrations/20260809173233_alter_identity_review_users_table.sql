-- +goose Up
-- +goose StatementBegin
ALTER TABLE identity_review_users
ADD identity_review_campaign_id INT NULL;

ALTER TABLE identity_review_users
ADD CONSTRAINT FK_identity_review_users_campaign
FOREIGN KEY (identity_review_campaign_id) REFERENCES identity_review_campaigns(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE identity_review_users
DROP CONSTRAINT FK_identity_review_users_campaign;

ALTER TABLE identity_review_users
DROP COLUMN identity_review_campaign_id;
-- +goose StatementEnd
