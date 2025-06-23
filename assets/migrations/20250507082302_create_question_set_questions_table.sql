-- +goose Up
-- +goose StatementBegin
CREATE TABLE question_set_questions (
    id INT PRIMARY KEY IDENTITY(1,1),
    question_set_id INT NOT NULL,
    third_party_question_library_id INT NOT NULL,
    question_title_en NVARCHAR(1000) NOT NULL,
    question_title_ar NVARCHAR(1000) NOT NULL,
    created_by_id INT,
    updated_by_id INT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,

    CONSTRAINT FK_qs_question_set_questions_set FOREIGN KEY (question_set_id) REFERENCES question_sets(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_qs_question_set_questions_library FOREIGN KEY (third_party_question_library_id) REFERENCES third_party_question_libraries(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_qs_question_set_questions_created_by FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_qs_question_set_questions_updated_by FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.question_set_questions;
-- +goose StatementEnd