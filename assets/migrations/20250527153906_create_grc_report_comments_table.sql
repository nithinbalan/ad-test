-- +goose Up
-- +goose StatementBegin
CREATE TABLE grc_report_workflow_review_comments (
    id INT PRIMARY KEY IDENTITY(1,1),
    grc_report_workflow_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    grc_report_comment_status_id INT NOT NULL,  
    resolved_by_id INT,
    grc_report_workflow_status_id INT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,

    CONSTRAINT FK_grc_report_workflow_review_comments_grc_report_workflows FOREIGN KEY (grc_report_workflow_id) REFERENCES grc_report_workflows(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_workflow_review_comments_grc_report_comment_statuses FOREIGN KEY (grc_report_comment_status_id) REFERENCES grc_report_comment_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_workflow_review_comments_grc_report_workflow_statuses FOREIGN KEY (grc_report_workflow_status_id) REFERENCES grc_report_workflow_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_workflow_review_comments_user_1 FOREIGN KEY (resolved_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_workflow_review_comments_user_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,  
    CONSTRAINT FK_grc_report_workflow_review_comments_user_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.grc_report_workflow_review_comments;
-- +goose StatementEnd
