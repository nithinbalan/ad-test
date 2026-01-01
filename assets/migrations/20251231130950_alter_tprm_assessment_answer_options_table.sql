-- +goose Up
-- +goose StatementBegin
ALTER TABLE tprm_assessment_answer_options
ADD tprm_compliance_response_id INT NULL;

ALTER TABLE tprm_assessment_answer_options
ADD CONSTRAINT FK_tprm_assessment_answer_options_tprm_compliance_responses
FOREIGN KEY (tprm_compliance_response_id) REFERENCES tprm_compliance_responses(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE tprm_assessment_answer_options
DROP CONSTRAINT FK_tprm_assessment_answer_options_tprm_compliance_responses;


ALTER TABLE tprm_assessment_answer_options
DROP COLUMN tprm_compliance_response_id;
-- +goose StatementEnd
