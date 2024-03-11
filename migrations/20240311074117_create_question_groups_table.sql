-- +goose Up
-- +goose StatementBegin
CREATE TABLE question_groups (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en VARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,    
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_question_groups_status CHECK (status IN ('active', 'inactive'))
);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.question_groups;
-- +goose StatementEnd
