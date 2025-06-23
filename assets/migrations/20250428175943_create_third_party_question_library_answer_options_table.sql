-- +goose Up
-- +goose StatementBegin
CREATE TABLE third_party_question_library_answer_options
(
    id INT PRIMARY KEY IDENTITY(1,1),
    third_party_question_library_id INT NOT NULL,
    answer_option_title_en NVARCHAR(500),
    answer_option_title_ar NVARCHAR(500),
    risk_level_id INT NULL,
    range_from INT NULL,
    range_to INT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NULL,

    CONSTRAINT FK_third_party_question_library_answer_options_third_party_question_library FOREIGN KEY (third_party_question_library_id) REFERENCES third_party_question_libraries(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_question_library_answer_options_risk_level FOREIGN KEY (risk_level_id) REFERENCES risk_levels(id) ON DELETE SET NULL ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_question_library_answer_options_created_by FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_third_party_question_library_answer_options_updated_by FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.third_party_question_library_answer_options;
-- +goose StatementEnd
