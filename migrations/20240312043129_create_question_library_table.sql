-- +goose Up
-- +goose StatementBegin
CREATE TABLE question_library (
    id INT PRIMARY KEY IDENTITY(1, 1),
    question_group_id INT NOT NULL,
    question_en VARCHAR(500) NOT NULL,
    question_ar NVARCHAR(500) NOT NULL,    
    CONSTRAINT FK_question_group_question_library FOREIGN KEY (question_group_id) REFERENCES question_groups(id),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_question_library_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.question_library;
-- +goose StatementEnd
