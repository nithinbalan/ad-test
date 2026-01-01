-- +goose Up
-- +goose StatementBegin

ALTER TABLE tasks
ADD tprm_assessment_section_question_action_plan_id INT NULL;

ALTER TABLE tasks
ADD CONSTRAINT FK_tasks_tprm_assessment_section_question_action_plans
    FOREIGN KEY (tprm_assessment_section_question_action_plan_id)
    REFERENCES tprm_assessment_section_question_action_plans(id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE tasks
DROP CONSTRAINT FK_tasks_tprm_assessment_section_question_action_plans;

ALTER TABLE tasks
DROP COLUMN tprm_assessment_section_question_action_plan_id;

-- +goose StatementEnd
