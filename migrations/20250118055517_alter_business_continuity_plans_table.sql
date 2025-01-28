-- +goose Up
-- +goose StatementBegin
ALTER TABLE business_continuity_plans 
ADD responsible_user_id INT NULL;

ALTER TABLE business_continuity_plans 
ADD responsible_unit_id INT NULL;

ALTER TABLE business_continuity_plans
ADD CONSTRAINT FK_business_continuity_plans_users_3 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE business_continuity_plans
ADD CONSTRAINT FK_business_continuity_plans_organization_units FOREIGN KEY (responsible_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
