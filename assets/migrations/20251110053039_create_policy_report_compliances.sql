-- +goose Up
-- +goose StatementBegin
CREATE TABLE policy_report_compliances (
    id INT PRIMARY KEY IDENTITY(1,1),
    policy_report_id INT NOT NULL,
    compliance_title_en NVARCHAR(255) NOT NULL,
    compliance_title_ar NVARCHAR(255) NULL,
    clauses NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_policy_report_compliances_reports
        FOREIGN KEY (policy_report_id) REFERENCES policy_reports (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_policy_report_compliances_users_1
        FOREIGN KEY (created_by_id) REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE policy_report_compliances;
-- +goose StatementEnd
