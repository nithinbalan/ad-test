-- +goose Up
-- +goose StatementBegin
ALTER TABLE bcp_recovery_plans
DROP CONSTRAINT FK_bcp_recovery_plans_time_units;

ALTER TABLE bcp_recovery_plans
ADD CONSTRAINT FK_bcp_recovery_plans_business_continuity_duration_units FOREIGN KEY (target_completion_unit_id) REFERENCES business_continuity_duration_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE bcp_recovery_plans
DROP CONSTRAINT FK_bcp_recovery_plans_business_continuity_duration_units;

ALTER TABLE bcp_recovery_plans
ADD CONSTRAINT FK_bcp_recovery_plans_time_units FOREIGN KEY (target_completion_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd
