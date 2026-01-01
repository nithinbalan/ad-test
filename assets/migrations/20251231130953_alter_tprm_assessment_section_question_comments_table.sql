-- +goose Up
-- +goose StatementBegin
ALTER TABLE tprm_assessment_section_question_comments 
ADD third_party_assessment_id INT NULL;

ALTER TABLE tprm_assessment_section_question_comments
ADD CONSTRAINT FK_tprm_assessment_section_question_comments_third_party_assessment
FOREIGN KEY (third_party_assessment_id) REFERENCES third_party_assessments(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE tprm_assessment_section_question_comments
ALTER COLUMN tprm_assessment_section_question_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tprm_assessment_section_question_comments
DROP CONSTRAINT FK_tprm_assessment_section_question_comments_third_party_assessment;

ALTER TABLE tprm_assessment_section_question_comments
DROP COLUMN third_party_assessment_id;

ALTER TABLE tprm_assessment_section_question_comments
ALTER COLUMN tprm_assessment_section_question_id INT NOT NULL;
-- +goose StatementEnd
