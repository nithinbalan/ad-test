-- +goose Up
-- +goose StatementBegin
CREATE TABLE awareness_reports(
    id INT PRIMARY KEY IDENTITY(1,1),
    reference_number VARCHAR(50) NOT NULL UNIQUE,
    awareness_report_type_id INT NOT NULL,
    sub_category_id INT NULL,
    report_source_channel VARCHAR(255) NULL,
    severity_id INT NULL,
    report_date DATETIME NULL,
    resulted_from_campaign BIT NOT NULL DEFAULT 0,
    awareness_handling_status_id INT NULL,
    handling_owner_id INT NULL,     
    closure_date DATETIME NULL,
    result NVARCHAR(MAX) NULL,

    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_awareness_reports_types FOREIGN KEY (awareness_report_type_id) REFERENCES awareness_campaign_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_reports_sub_categories FOREIGN KEY (sub_category_id) REFERENCES awareness_report_sub_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_reports_severity_levels FOREIGN KEY (severity_id) REFERENCES severity_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_reports_handling_statuses FOREIGN KEY (awareness_handling_status_id) REFERENCES awareness_handling_statuses(id)  ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_reports_users_1 FOREIGN KEY (handling_owner_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_reports_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_reports_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE awareness_reports;
-- +goose StatementEnd
