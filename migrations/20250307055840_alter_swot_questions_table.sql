-- +goose Up
-- +goose StatementBegin

EXEC sp_rename 'swot_questions', 'swot_pestel_questions';

ALTER TABLE swot_pestel_questions
ADD is_swot TINYINT NULL;

ALTER TABLE swot_pestel_questions
ADD is_pestel TINYINT NULL;

ALTER TABLE swot_pestel_questions
DROP CONSTRAINT FK_swot_questions_swot_question_groups;

EXEC sp_rename 'swot_pestel_questions.swot_question_group_id', 'swot_pestel_question_group_id', 'COLUMN';

ALTER TABLE swot_pestel_questions
ADD CONSTRAINT FK_swot_pestel_questions_swot_pestel_question_groups
FOREIGN KEY (swot_pestel_question_group_id) REFERENCES swot_pestel_question_groups(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE swot_pestel_questions
DROP COLUMN is_swot;

ALTER TABLE swot_pestel_questions
DROP COLUMN is_pestel;

ALTER TABLE swot_pestel_questions
DROP CONSTRAINT FK_swot_questions_swot_question_groups;

ALTER TABLE swot_pestel_questions
ADD CONSTRAINT FK_swot_questions_swot_question_groups
FOREIGN KEY (swot_question_group_id) REFERENCES swot_question_groups(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

EXEC sp_rename 'swot_pestel_questions', 'swot_questions';
-- +goose StatementEnd
