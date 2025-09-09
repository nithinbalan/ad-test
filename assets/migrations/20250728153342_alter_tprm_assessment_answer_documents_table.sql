-- +goose Up
-- +goose StatementBegin
ALTER TABLE tprm_assessment_answer_documents
ADD tprm_assessment_answer_id INT NULL;

ALTER TABLE tprm_assessment_answer_documents
ADD CONSTRAINT FK_tprm_assessment_answer_documents_answers 
    FOREIGN KEY (tprm_assessment_answer_id) REFERENCES tprm_assessment_answers (id) 
    ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tprm_assessment_answer_documents
DROP CONSTRAINT FK_tprm_assessment_answer_documents_answers;


ALTER TABLE tprm_assessment_answer_documents
DROP COLUMN tprm_assessment_answer_id;
-- +goose StatementEnd
