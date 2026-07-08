-- +goose Up
-- +goose StatementBegin
CREATE TABLE bia_process_workflow_reviews (
    id INT PRIMARY KEY IDENTITY(1,1),
    bia_process_workflow_id INT NOT NULL,
    workflow_status_id INT NOT NULL,
    comment NVARCHAR(MAX),
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_bia_process_wrv_reviews_workflow
        FOREIGN KEY (bia_process_workflow_id)
        REFERENCES bia_process_workflows(id),

    CONSTRAINT FK_bia_process_wrv_reviews_workflow_status
        FOREIGN KEY (workflow_status_id)
        REFERENCES workflow_statuses(id),

    CONSTRAINT FK_bia_process_wrv_reviews_users_1
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_bia_process_wrv_reviews_users_2
        FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE bia_process_workflow_reviews;
-- +goose StatementEnd
