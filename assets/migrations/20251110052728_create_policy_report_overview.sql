-- +goose Up
-- +goose StatementBegin
CREATE TABLE policy_report_overviews (
    id INT PRIMARY KEY IDENTITY(1,1),
    policy_report_id INT NOT NULL,
    reference_number VARCHAR(50) NULL,
    policy_title_en NVARCHAR(255) NOT NULL,
    policy_title_ar NVARCHAR(255) NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    responsible_unit_title_en NVARCHAR(255) NULL,
    responsible_unit_title_ar NVARCHAR(255) NULL,
    responsible_user_name_en NVARCHAR(255) NULL,
    responsible_user_name_ar NVARCHAR(255) NULL,
    published_date VARCHAR(50) NULL,
    next_review_date VARCHAR(50) NULL,
    document_version DECIMAL(4,2) NULL,
    document_status_title_en NVARCHAR(100) NULL,
    document_status_title_ar NVARCHAR(100) NULL,
    document_status_color_code NVARCHAR(100) NULL,
    document_access_type_title_en NVARCHAR(100) NULL,
    document_access_type_title_ar NVARCHAR(100) NULL,
    document_access_type_color_code VARCHAR(9) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_policy_report_overviews_reports 
        FOREIGN KEY (policy_report_id) REFERENCES policy_reports (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_policy_report_overviews_users_1 
        FOREIGN KEY (created_by_id) REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE policy_report_overviews;
-- +goose StatementEnd
