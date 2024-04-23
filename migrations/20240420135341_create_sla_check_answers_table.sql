-- +goose Up
-- +goose StatementBegin
CREATE TABLE sla_check_answers (
    id INT PRIMARY KEY IDENTITY(1, 1),
    sla_check_question_id INT NOT NULL,
    compliance_status_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    CONSTRAINT FK_sla_check_answers_sla_check_questions FOREIGN KEY (sla_check_question_id) REFERENCES sla_check_questions (id),
    CONSTRAINT FK_sla_check_answers_compliance_statuses FOREIGN KEY (compliance_status_id) REFERENCES compliance_statuses (id),
    CONSTRAINT FK_sla_check_answers_users_1 FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT FK_sla_check_answers_users_2 FOREIGN KEY (updated_by) REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sla_check_answers;
-- +goose StatementEnd
