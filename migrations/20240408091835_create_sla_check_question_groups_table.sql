-- +goose Up
-- +goose StatementBegin
CREATE TABLE sla_check_question_groups (
    id INT PRIMARY KEY IDENTITY(1, 1),
    sla_check_id INT NOT NULL,
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    total_count INT NOT NULL DEFAULT 0,
    completed_count INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    CONSTRAINT FK_sla_check_question_groups_sla_checks FOREIGN KEY (sla_check_id) REFERENCES sla_checks (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_check_question_groups_users_1 FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_check_question_groups_users_2 FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sla_check_question_groups;
-- +goose StatementEnd
