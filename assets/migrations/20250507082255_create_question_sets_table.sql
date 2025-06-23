-- +goose Up
-- +goose StatementBegin
CREATE TABLE question_sets (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX),
    created_by_id INT,
    updated_by_id INT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,

    CONSTRAINT FK_qs_question_sets_created_by FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_qs_question_sets_updated_by FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.question_sets;
-- +goose StatementEnd