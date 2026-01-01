-- +goose Up
-- +goose StatementBegin
ALTER TABLE third_party_question_library_answer_options 
ADD tprm_compliance_response_id INT NULL;

ALTER TABLE third_party_question_library_answer_options
ADD CONSTRAINT FK_third_party_question_library_answer_options_responses
FOREIGN KEY (tprm_compliance_response_id) REFERENCES tprm_compliance_responses(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE third_party_question_library_answer_options
DROP CONSTRAINT FK_third_party_question_library_answer_options_responses;

ALTER TABLE third_party_question_library_answer_options
DROP COLUMN tprm_compliance_response_id;
-- +goose StatementEnd
