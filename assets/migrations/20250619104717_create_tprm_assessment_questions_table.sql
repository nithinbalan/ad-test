-- +goose Up
-- +goose StatementBegin
CREATE TABLE tprm_assessment_questions (
    id INT PRIMARY KEY IDENTITY(1,1),
    third_party_assessment_id INT NOT NULL,
    third_party_question_library_id INT NOT NULL,
    question_title_en NVARCHAR(1000) NOT NULL,
    question_title_ar NVARCHAR(1000) NOT NULL,
    created_by_id INT,
    updated_by_id INT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,

    CONSTRAINT FK_tprm_assessment_questions_third_party_assessments FOREIGN KEY (third_party_assessment_id) REFERENCES third_party_assessments(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tprm_assessment_questions_library FOREIGN KEY (third_party_question_library_id) REFERENCES third_party_question_libraries(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tprm_assessment_questions_created_by FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tprm_assessment_questions_updated_by FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.tprm_assessment_questions;
-- +goose StatementEnd