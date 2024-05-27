-- +goose Up
-- +goose StatementBegin
CREATE TABLE sla_check_answer_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    sla_check_answer_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    type VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    icon VARCHAR(10) NOT NULL,
    created_by INT,
    updated_by INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_sla_check_answer_documents FOREIGN KEY (sla_check_answer_id) REFERENCES sla_check_answers(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_check_answer_documents_users_1 FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_check_answer_documents_users_2 FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.sla_check_answer_documents;
-- +goose StatementEnd
