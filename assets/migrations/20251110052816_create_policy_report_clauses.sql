-- +goose Up
-- +goose StatementBegin
CREATE TABLE policy_report_clauses (
    id INT PRIMARY KEY IDENTITY(1,1),
    policy_report_id INT NOT NULL,
    clause_number VARCHAR(50) NULL,
    clause_title_en NVARCHAR(255) NOT NULL,
    clause_title_ar NVARCHAR(255) NULL,
    description NVARCHAR(MAX) NULL,
    controls NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_policy_report_clauses_reports
        FOREIGN KEY (policy_report_id) REFERENCES policy_reports (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_policy_report_clauses_users_1
        FOREIGN KEY (created_by_id) REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE policy_report_clauses;
-- +goose StatementEnd
