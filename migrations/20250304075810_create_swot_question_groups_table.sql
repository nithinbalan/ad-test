-- +goose Up
-- +goose StatementBegin
CREATE TABLE swot_question_groups (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(MAX) NOT NULL,
    title_ar NVARCHAR(MAX) NOT NULL,
    [order] INT,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_swot_question_groups_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE swot_question_groups;
-- +goose StatementEnd
