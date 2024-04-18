-- +goose Up
-- +goose StatementBegin
CREATE TABLE sla_check_questions (
    id INT PRIMARY KEY IDENTITY(1, 1),
    sla_check_question_group_id INT NOT NULL,
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    CONSTRAINT FK_sla_check_questions_sla_check_question_groups FOREIGN KEY (sla_check_question_group_id) REFERENCES sla_check_question_groups (id),
    CONSTRAINT FK_sla_check_questions_users1 FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT FK_sla_check_questions_users2 FOREIGN KEY (updated_by) REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sla_check_questions;
-- +goose StatementEnd
