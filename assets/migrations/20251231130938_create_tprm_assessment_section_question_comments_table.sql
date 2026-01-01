-- +goose Up
-- +goose StatementBegin
CREATE TABLE tprm_assessment_section_question_comments
(
    id INT PRIMARY KEY IDENTITY(1,1),
    tprm_assessment_section_question_id INT NOT NULL,
    tprm_assessment_section_question_comment_id INT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    level INT NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_tprm_assessment_section_question_comments_question
        FOREIGN KEY (tprm_assessment_section_question_id)
        REFERENCES tprm_assessment_section_questions(id)
        ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT FK_tprm_assessment_section_question_comments_parent
        FOREIGN KEY (tprm_assessment_section_question_comment_id)
        REFERENCES tprm_assessment_section_question_comments(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_tprm_assessment_section_question_comments_users_1
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_tprm_assessment_section_question_comments_users_2
        FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tprm_assessment_section_question_comments;
-- +goose StatementEnd
