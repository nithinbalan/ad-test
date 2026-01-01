-- +goose Up
-- +goose StatementBegin
CREATE TABLE question_set_sections (
    id INT IDENTITY(1,1) PRIMARY KEY,
    question_set_id INT NOT NULL,
    third_party_question_category_id INT NULL,
    title_en NVARCHAR(MAX) NOT NULL,
    title_ar NVARCHAR(MAX) NOT NULL,
    [order] INT NOT NULL,  -- order is a reserved keyword
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),

    CONSTRAINT FK_question_set_sections_question_set_id
        FOREIGN KEY (question_set_id)
        REFERENCES question_sets(id)
        ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT FK_question_set_sections_third_party_question_category
        FOREIGN KEY (third_party_question_category_id)
        REFERENCES third_party_question_categories(id)
        ON DELETE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE question_set_sections;
-- +goose StatementEnd
