-- +goose Up
-- +goose StatementBegin
CREATE TABLE internal_committee_workflow_review_history (
    id INT PRIMARY KEY IDENTITY(1,1),
    internal_committee_workflow_id INT NOT NULL,
    workflow_status_id INT NOT NULL,
    comment NVARCHAR(MAX),
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,

    CONSTRAINT FK_internal_committee_workflow_review_history_workflow
        FOREIGN KEY (internal_committee_workflow_id)
        REFERENCES internal_committee_workflows(id)
       ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_internal_committee_workflow_review_history_status
        FOREIGN KEY (workflow_status_id)
        REFERENCES workflow_statuses(id)
       ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_internal_committee_workflow_review_history_users_1
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
       ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_internal_committee_workflow_review_history_users_2
        FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
       ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE internal_committee_workflow_review_history;
-- +goose StatementEnd