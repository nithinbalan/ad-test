-- +goose Up
-- +goose StatementBegin
CREATE TABLE tprm_assessment_answer_comment_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    tprm_assessment_question_id INT NOT NULL,
    tprm_assessment_answer_id INT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NULL,
    CONSTRAINT FK_tprm_assessment_answer_comment_document_questions FOREIGN KEY (tprm_assessment_question_id) REFERENCES tprm_assessment_questions(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tprm_assessment_answer_comment_document_answers FOREIGN KEY (tprm_assessment_answer_id) REFERENCES tprm_assessment_answers(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.tprm_assessment_answer_comment_documents;
-- +goose StatementEnd
