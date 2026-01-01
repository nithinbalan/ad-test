-- +goose Up
-- +goose StatementBegin
CREATE TABLE tprm_assessment_sections (
    id INT IDENTITY(1,1) PRIMARY KEY,
    third_party_assessment_id INT NOT NULL,
    third_party_question_category_id INT NULL,
    title_en NVARCHAR(MAX) NOT NULL,
    title_ar NVARCHAR(MAX) NOT NULL,
    [order] INT NOT NULL,  -- order is a reserved keyword
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),

    CONSTRAINT FK_tprm_assessment_sections_third_party_assessment
        FOREIGN KEY (third_party_assessment_id)
        REFERENCES third_party_assessments(id)
        ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT FK_tprm_assessment_sections_third_party_question_category
        FOREIGN KEY (third_party_question_category_id)
        REFERENCES third_party_question_categories(id)
        ON DELETE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tprm_assessment_sections;
-- +goose StatementEnd
