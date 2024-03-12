-- +goose Up
-- +goose StatementBegin
CREATE TABLE question_library (
    id INT PRIMARY KEY IDENTITY(1, 1),
    question_group_id INT NOT NULL,
    question_en VARCHAR(MAX) NOT NULL,
    question_ar NVARCHAR(MAX) NOT NULL,    
    CONSTRAINT FK_question_group_question_library FOREIGN KEY (question_group_id) REFERENCES question_groups(id),
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.question_library;
-- +goose StatementEnd
