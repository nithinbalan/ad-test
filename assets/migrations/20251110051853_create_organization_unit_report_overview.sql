-- +goose Up
-- +goose StatementBegin
CREATE TABLE organization_unit_report_overviews (
    id INT PRIMARY KEY IDENTITY(1,1),
    organization_unit_report_id INT NOT NULL,
    organization_unit_title_en NVARCHAR(255) NOT NULL,
    organization_unit_title_ar NVARCHAR(255) NULL,
    organization_unit_head_name_en NVARCHAR(255) NULL,
    organization_unit_head_name_ar NVARCHAR(255) NULL,
    organization_unit_head_email VARCHAR(255) NULL,
    total_employees INT NOT NULL DEFAULT 0,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    
    CONSTRAINT FK_organization_unit_report_overviews_reports 
        FOREIGN KEY (organization_unit_report_id) REFERENCES organization_unit_reports (id) 
        ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_organization_unit_report_overviews_users_1 
        FOREIGN KEY (created_by_id) REFERENCES users (id) 
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE organization_unit_report_overviews;
-- +goose StatementEnd
