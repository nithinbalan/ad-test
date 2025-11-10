-- +goose Up
-- +goose StatementBegin
CREATE TABLE organization_unit_report_policies (
    id INT PRIMARY KEY IDENTITY(1,1),
    organization_unit_report_id INT NOT NULL,
    document_title_en NVARCHAR(255) NOT NULL,
    document_title_ar NVARCHAR(255),
    document_category_title_en NVARCHAR(255),
    document_category_title_ar NVARCHAR(255),
    document_type_title_en NVARCHAR(100),
    document_type_title_ar NVARCHAR(100),
    responsible_user_name_en NVARCHAR(255),
    responsible_user_name_ar NVARCHAR(255),
    acknowledged BIT NOT NULL DEFAULT 0,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_organization_unit_report_policies_reports 
        FOREIGN KEY (organization_unit_report_id) REFERENCES organization_unit_reports (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_organization_unit_report_policies_users_1 
        FOREIGN KEY (created_by_id) REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE organization_unit_report_policies;
-- +goose StatementEnd
