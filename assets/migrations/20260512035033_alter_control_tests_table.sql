-- +goose Up
-- +goose StatementBegin
ALTER TABLE control_tests
ADD
    risk_register_update_plan_id INT NULL,

    CONSTRAINT FK_control_test_risk_register_update_plan
        FOREIGN KEY (risk_register_update_plan_id)
        REFERENCES risk_register_update_plans(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE control_tests
DROP CONSTRAINT FK_control_test_risk_register_update_plan;

ALTER TABLE control_tests
DROP COLUMN risk_register_update_plan_id;
-- +goose StatementEnd
