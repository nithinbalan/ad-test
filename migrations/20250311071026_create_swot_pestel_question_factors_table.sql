-- +goose Up
-- +goose StatementBegin
CREATE TABLE swot_pestel_question_factors (
    id INT PRIMARY KEY IDENTITY(1, 1),
    swot_pestel_question_id INT NOT NULL,
    swot_pestel_category_id INT NOT NULL,
    description NVARCHAR(MAX) NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_swot_pestel_question_factors_swot_pestel_questions FOREIGN KEY (swot_pestel_question_id) REFERENCES swot_pestel_questions(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_swot_pestel_question_factors_swot_pestel_categories FOREIGN KEY (swot_pestel_category_id) REFERENCES swot_pestel_categories(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_swot_pestel_question_factors_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id),
    CONSTRAINT FK_swot_pestel_question_factors_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id)

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.swot_pestel_question_factors;
-- +goose StatementEnd
