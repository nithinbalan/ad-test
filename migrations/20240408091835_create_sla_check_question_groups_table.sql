-- +goose Up
-- +goose StatementBegin
CREATE TABLE sla_check_question_groups (
    id INT PRIMARY KEY IDENTITY(1, 1),
    sla_check_id INT NOT NULL,
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    CONSTRAINT FK_sla_check_question_groups_sla_checks FOREIGN KEY (sla_check_id) REFERENCES sla_checks (id),
    CONSTRAINT FK_sla_check_question_groups_users1 FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT FK_sla_check_question_groups_users2 FOREIGN KEY (updated_by) REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sla_check_question_groups;
-- +goose StatementEnd
