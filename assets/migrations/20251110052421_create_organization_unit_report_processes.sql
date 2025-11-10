-- +goose Up
-- +goose StatementBegin
CREATE TABLE organization_unit_report_processes (
    id INT PRIMARY KEY IDENTITY(1,1),
    organization_unit_report_id INT NOT NULL,
    process_title_en NVARCHAR(255) NOT NULL,
    process_title_ar NVARCHAR(255),
    process_group_title_en NVARCHAR(255),
    process_group_title_ar NVARCHAR(255),
    process_owner_name_en NVARCHAR(255),
    process_owner_name_ar NVARCHAR(255),
    responsible_unit_title_en NVARCHAR(255),
    responsible_unit_title_ar NVARCHAR(255),
    risk_rating_title_en NVARCHAR(100),
    risk_rating_title_ar NVARCHAR(100),
    risk_rating_color_code VARCHAR(9) NULL,
    lifecycle_status_title_en NVARCHAR(100),
    lifecycle_status_title_ar NVARCHAR(100),
    lifecycle_status_color_code VARCHAR(9) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_organization_unit_report_processes_reports 
        FOREIGN KEY (organization_unit_report_id) REFERENCES organization_unit_reports (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_organization_unit_report_processes_users_1 
        FOREIGN KEY (created_by_id) REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE organization_unit_report_processes;
-- +goose StatementEnd
