-- +goose Up
-- +goose StatementBegin
ALTER TABLE identity_review_campaigns
DROP CONSTRAINT FK_identity_review_campaigns_users_2;

ALTER TABLE identity_review_campaigns
DROP COLUMN approver_id;

ALTER TABLE identity_review_campaigns
ADD workflow_id INT NOT NULL;

ALTER TABLE identity_review_campaigns
ADD CONSTRAINT FK_identity_review_campaigns_workflows
FOREIGN KEY (workflow_id) REFERENCES workflows(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE identity_review_campaigns
DROP CONSTRAINT FK_identity_review_campaigns_workflows;

ALTER TABLE identity_review_campaigns
DROP COLUMN workflow_id;

ALTER TABLE identity_review_campaigns
ADD approver_id INT NULL;

ALTER TABLE identity_review_campaigns
ADD CONSTRAINT FK_identity_review_campaigns_users_2
FOREIGN KEY (approver_id) REFERENCES users(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;
-- +goose StatementEnd
