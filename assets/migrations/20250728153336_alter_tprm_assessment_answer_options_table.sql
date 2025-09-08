-- +goose Up
-- +goose StatementBegin
ALTER TABLE tprm_assessment_answer_options
ADD tprm_assessment_answer_id INT NULL;

ALTER TABLE tprm_assessment_answer_options
ADD CONSTRAINT FK_tprm_assessment_answer_options_answers 
    FOREIGN KEY (tprm_assessment_answer_id) REFERENCES tprm_assessment_answers (id) 
    ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tprm_assessment_answer_options
DROP CONSTRAINT FK_tprm_assessment_answer_options_answers;


ALTER TABLE tprm_assessment_answer_options
DROP COLUMN tprm_assessment_answer_id;
-- +goose StatementEnd
