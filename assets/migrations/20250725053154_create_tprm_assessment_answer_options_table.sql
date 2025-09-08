-- +goose Up
-- +goose StatementBegin
CREATE TABLE tprm_assessment_answer_options (
    id INT PRIMARY KEY IDENTITY(1,1),
    tprm_assessment_question_id INT NOT NULL,
    answer_option_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,

    CONSTRAINT FK_tprm_assessment_answer_options_tprm_assessment_questions FOREIGN KEY (tprm_assessment_question_id) REFERENCES tprm_assessment_questions(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tprm_assessment_answer_options_tprm_assessment_question_answer_options FOREIGN KEY (answer_option_id) REFERENCES tprm_assessment_question_answer_options(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.tprm_assessment_answer_options;
-- +goose StatementEnd