-- +goose Up
-- +goose StatementBegin
-- Rename the table
EXEC sp_rename 'tprm_assessment_questions', 'tprm_assessment_section_questions';

-- Rename the constraints
EXEC sp_rename 'FK_tprm_assessment_questions_third_party_assessments', 'FK_tprm_assessment_section_questions_third_party_assessments', 'OBJECT';
EXEC sp_rename 'FK_tprm_assessment_questions_library', 'FK_tprm_assessment_section_questions_library', 'OBJECT';
EXEC sp_rename 'FK_tprm_assessment_questions_created_by', 'FK_tprm_assessment_section_questions_created_by', 'OBJECT';
EXEC sp_rename 'FK_tprm_assessment_questions_updated_by', 'FK_tprm_assessment_section_questions_updated_by', 'OBJECT';

ALTER TABLE tprm_assessment_section_questions
    ADD [order] INT NULL;

ALTER TABLE tprm_assessment_section_questions
    ADD tprm_assessment_section_id INT NULL;

ALTER TABLE tprm_assessment_section_questions
ADD CONSTRAINT FK_tprm_assessment_section_questions_tprm_assessment_section
FOREIGN KEY (tprm_assessment_section_id) REFERENCES tprm_assessment_sections(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Rename the table back to its original name

ALTER TABLE tprm_assessment_section_questions
DROP CONSTRAINT FK_tprm_assessment_section_questions_tprm_assessment_section;

ALTER TABLE tprm_assessment_section_questions
DROP COLUMN tprm_assessment_section_id;

ALTER TABLE tprm_assessment_section_questions
DROP COLUMN [order];

EXEC sp_rename 'tprm_assessment_section_questions', 'tprm_assessment_questions';

-- Rename the constraints back to their original names
EXEC sp_rename 'FK_tprm_assessment_section_questions_third_party_assessments', 'FK_tprm_assessment_questions_third_party_assessments', 'OBJECT';
EXEC sp_rename 'FK_tprm_assessment_section_questions_library', 'FK_tprm_assessment_questions_library', 'OBJECT';
EXEC sp_rename 'FK_tprm_assessment_section_questions_created_by', 'FK_tprm_assessment_questions_created_by', 'OBJECT';
EXEC sp_rename 'FK_tprm_assessment_section_questions_updated_by', 'FK_tprm_assessment_questions_updated_by', 'OBJECT';
-- +goose StatementEnd
