-- +goose Up
-- +goose StatementBegin
-- Rename the table
EXEC sp_rename 'question_set_questions', 'question_set_section_questions';

-- Rename the constraints
EXEC sp_rename 'FK_qs_question_set_questions_set', 'FK_question_set_section_questions_questions_set', 'OBJECT';
EXEC sp_rename 'FK_qs_question_set_questions_library', 'FK_question_set_section_questions_questions_library', 'OBJECT';
EXEC sp_rename 'FK_qs_question_set_questions_created_by', 'FK_question_set_section_questions_created_by', 'OBJECT';
EXEC sp_rename 'FK_qs_question_set_questions_updated_by', 'FK_question_set_section_questions_updated_by', 'OBJECT';

ALTER TABLE question_set_section_questions
    ADD [order] INT NULL;

ALTER TABLE question_set_section_questions
    ADD question_set_section_id INT NULL;

ALTER TABLE question_set_section_questions
ADD CONSTRAINT FK_question_set_section_questions_question_set_section
FOREIGN KEY (question_set_section_id) REFERENCES question_set_sections(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Rename the table back to its original name

ALTER TABLE question_set_section_questions
DROP CONSTRAINT FK_question_set_section_questions_question_set_section;

ALTER TABLE question_set_section_questions
DROP COLUMN question_set_section_id;

ALTER TABLE question_set_section_questions
DROP COLUMN [order];

EXEC sp_rename 'question_set_section_questions', 'question_set_questions';

-- Rename the constraints back to their original names
EXEC sp_rename 'FK_question_set_section_questions_questions_set', 'FK_qs_question_set_questions_set', 'OBJECT';
EXEC sp_rename 'FK_question_set_section_questions_questions_library', 'FK_qs_question_set_questions_library', 'OBJECT';
EXEC sp_rename 'FK_question_set_section_questions_created_by', 'FK_qs_question_set_questions_created_by', 'OBJECT';
EXEC sp_rename 'FK_question_set_section_questions_updated_by', 'FK_qs_question_set_questions_updated_by', 'OBJECT';
-- +goose StatementEnd
