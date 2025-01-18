-- +goose Up
-- +goose StatementBegin
ALTER TABLE business_continuity_plans
ADD bcp_status_id INT NULL;

ALTER TABLE business_continuity_plans
ADD CONSTRAINT FK_business_continuity_plans_bcp_statuses FOREIGN KEY (bcp_status_id) REFERENCES bcp_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE business_continuity_plans
DROP CONSTRAINT FK_business_continuity_plans_bcp_statuses;

ALTER TABLE business_continuity_plans
DROP COLUMN bcp_status_id;
-- +goose StatementEnd
