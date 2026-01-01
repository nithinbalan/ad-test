-- +goose Up
-- +goose StatementBegin

ALTER TABLE tprm_assessment_section_question_action_plans
ADD delegatee_user_id INT NULL;

ALTER TABLE tprm_assessment_section_question_action_plans
ADD CONSTRAINT FK_tprm_asq_action_plans_delegatee_user
    FOREIGN KEY (delegatee_user_id) REFERENCES users(id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE tprm_assessment_section_question_action_plans
DROP CONSTRAINT FK_tprm_asq_action_plans_delegatee_user;

ALTER TABLE tprm_assessment_section_question_action_plans
DROP COLUMN delegatee_user_id;

-- +goose StatementEnd
