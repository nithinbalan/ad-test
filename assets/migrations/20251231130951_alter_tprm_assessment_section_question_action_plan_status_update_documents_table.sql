-- +goose Up
-- +goose StatementBegin
ALTER TABLE tprm_assessment_section_question_action_plan_status_update_documents
ADD original_name NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tprm_assessment_section_question_action_plan_status_update_documents
DROP COLUMN original_name;
-- +goose StatementEnd
