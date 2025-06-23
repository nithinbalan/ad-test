-- +goose Up
-- +goose StatementBegin
CREATE TABLE grc_report_workflows (
    id INT PRIMARY KEY IDENTITY(1,1),
    grc_report_id INT NOT NULL,
    type VARCHAR(100) NOT NULL,
     CONSTRAINT CHK_grc_report_workflows_type CHECK (type IN ('governance','risk','compliance')),
    grc_report_workflow_status_id INT NOT NULL,  
    review_user_id INT NULL,
    approve_user_id INT NULL,
    reviewed_at DATETIME NULL,
    approved_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,

    CONSTRAINT FK_grc_report_workflows_grc_reports FOREIGN KEY (grc_report_id) REFERENCES grc_reports(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_workflows_grc_report_workflow_statuses FOREIGN KEY (grc_report_workflow_status_id) REFERENCES grc_report_workflow_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_workflows_user_1 FOREIGN KEY (review_user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_workflows_user_2 FOREIGN KEY (approve_user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_workflows_user_3 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,  
    CONSTRAINT FK_grc_report_workflows_user_4 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.grc_report_workflows;
-- +goose StatementEnd
