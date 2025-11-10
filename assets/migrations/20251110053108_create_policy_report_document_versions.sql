-- +goose Up
-- +goose StatementBegin
CREATE TABLE policy_report_document_versions (
    id INT PRIMARY KEY IDENTITY(1,1),
    policy_report_id INT NOT NULL,
    document_version DECIMAL(4,2) NOT NULL,
    published_at VARCHAR(50) NULL,
    published_user_name_en NVARCHAR(255) NULL,
    published_user_name_ar NVARCHAR(255) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_policy_report_document_versions_reports
        FOREIGN KEY (policy_report_id) REFERENCES policy_reports (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_policy_report_document_versions_users_1
        FOREIGN KEY (created_by_id) REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE policy_report_document_versions;
-- +goose StatementEnd
