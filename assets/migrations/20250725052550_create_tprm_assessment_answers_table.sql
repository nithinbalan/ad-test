-- +goose Up
-- +goose StatementBegin
CREATE TABLE tprm_assessment_answers (
    id INT PRIMARY KEY IDENTITY(1,1),
    tprm_assessment_question_id INT NOT NULL,
    rating_scale INT NULL,
    answer_text NVARCHAR(500) NULL,
    answer_date DATE NULL,
    not_applicable TINYINT DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,

    CONSTRAINT FK_tprm_assessment_answers_tprm_assessment_questions FOREIGN KEY (tprm_assessment_question_id) REFERENCES tprm_assessment_questions(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.tprm_assessment_answers;
-- +goose StatementEnd