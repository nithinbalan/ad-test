-- +goose Up
-- +goose StatementBegin
CREATE TABLE swot_questions (
    id INT PRIMARY KEY IDENTITY(1, 1),
    swot_question_group_id INT NOT NULL,
    title_en NVARCHAR(MAX) NOT NULL,
    title_ar NVARCHAR(MAX) NOT NULL,
    [order] INT,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_swot_questions_status CHECK (status IN ('active', 'inactive')),
    CONSTRAINT FK_swot_questions_swot_question_groups FOREIGN KEY (swot_question_group_id) REFERENCES swot_question_groups (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE swot_questions;
-- +goose StatementEnd
