-- +goose Up
-- +goose StatementBegin
CREATE TABLE third_party_question_libraries
(
    id INT PRIMARY KEY IDENTITY(1,1),
    reference_number VARCHAR(50) NOT NULL UNIQUE,
    question_title_en NVARCHAR(1000) NOT NULL,
    question_title_ar NVARCHAR(1000) NOT NULL,
    third_party_question_category_id INT NOT NULL,
    mandatory BIT NOT NULL,
    third_party_question_type_id INT NOT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NULL,

    CONSTRAINT FK_third_party_question_libraries_category FOREIGN KEY (third_party_question_category_id) REFERENCES third_party_question_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_question_libraries_type FOREIGN KEY (third_party_question_type_id) REFERENCES third_party_question_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_question_libraries_created_by FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_question_libraries_updated_by FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.third_party_question_libraries;
-- +goose StatementEnd
