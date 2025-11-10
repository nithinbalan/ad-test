-- +goose Up
-- +goose StatementBegin
CREATE TABLE policy_report_document_approval_histories (
    id INT PRIMARY KEY IDENTITY(1,1),
    policy_report_id INT NOT NULL,
    approved_user_name_en NVARCHAR(255) NULL,
    approved_user_name_ar NVARCHAR(255) NULL,
    designation_title_en NVARCHAR(255) NULL,
    designation_title_ar NVARCHAR(255) NULL,
    action_title_en NVARCHAR(100) NULL,
    action_title_ar NVARCHAR(100) NULL,
    action_color_code VARCHAR(9) NULL,
    action_date VARCHAR(50) NULL,
    comment NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_policy_report_document_approval_histories_reports
        FOREIGN KEY (policy_report_id) REFERENCES policy_reports (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_policy_report_document_approval_histories_users_1
        FOREIGN KEY (created_by_id) REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE policy_report_document_approval_histories;
-- +goose StatementEnd
