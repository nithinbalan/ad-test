-- +goose Up
-- +goose StatementBegin

ALTER TABLE strategy_profile_action_plans
ADD description NVARCHAR(MAX) NULL;

ALTER TABLE strategy_profile_action_plans
ADD severity_level_id INT NULL;

ALTER TABLE strategy_profile_action_plans
ADD percentage INT NULL;

ALTER TABLE strategy_profile_action_plans
ADD comments NVARCHAR(MAX) NULL;

ALTER TABLE strategy_profile_action_plans
ADD responsible_unit_id INT NULL;

ALTER TABLE strategy_profile_action_plans
ADD review_user_id INT NULL;

ALTER TABLE strategy_profile_action_plans
ADD action_plan_status_id INT NULL;


ALTER TABLE strategy_profile_action_plans
ADD CONSTRAINT FK_strategy_profile_action_plans_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE strategy_profile_action_plans
ADD CONSTRAINT FK_strategy_profile_action_plans_users_4 FOREIGN KEY (review_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE strategy_profile_action_plans
ADD CONSTRAINT FK_strategy_profile_action_plans_organization_units FOREIGN KEY (responsible_unit_id) REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE strategy_profile_action_plans
ADD CONSTRAINT FK_strategy_profile_action_plans_action_plan_statuses FOREIGN KEY (action_plan_status_id) REFERENCES action_plan_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE strategy_profile_action_plans
DROP CONSTRAINT FK_strategy_profile_action_plans_severity_levels;

ALTER TABLE strategy_profile_action_plans
DROP COLUMN severity_level_id;

ALTER TABLE strategy_profile_action_plans
DROP CONSTRAINT FK_strategy_profile_action_plans_action_plan_statuses;

ALTER TABLE strategy_profile_action_plans
DROP COLUMN severity_level_id;

ALTER TABLE strategy_profile_action_plans
DROP CONSTRAINT FK_strategy_profile_action_plans_users_4;

ALTER TABLE strategy_profile_action_plans
DROP COLUMN review_user_id;

ALTER TABLE strategy_profile_action_plans
DROP CONSTRAINT FK_strategy_profile_action_plans_organization_units;

ALTER TABLE strategy_profile_action_plans
DROP COLUMN responsible_unit_id;

ALTER TABLE strategy_profile_action_plans
DROP CONSTRAINT FK_strategy_profile_action_plans_action_plan_statuses;

ALTER TABLE strategy_profile_action_plans
DROP COLUMN action_plan_status_id;

ALTER TABLE strategy_profile_action_plans
DROP COLUMN description;

ALTER TABLE strategy_profile_action_plans
DROP COLUMN percentage;

ALTER TABLE strategy_profile_action_plans
DROP COLUMN comments;

-- +goose StatementEnd
