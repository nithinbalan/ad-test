-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_response_workflow_reviews (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_response_workflow_id INT NOT NULL,
    workflow_status_id INT NOT NULL,
    comment NVARCHAR(MAX),
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_rrw_reviews_workflow
        FOREIGN KEY (risk_response_workflow_id)
        REFERENCES risk_response_workflows(id),

    CONSTRAINT FK_rrw_reviews_workflow_status
        FOREIGN KEY (workflow_status_id)
        REFERENCES workflow_statuses(id),

    CONSTRAINT FK_rrw_reviews_users_1
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_rrw_reviews_users_2
        FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.risk_response_workflow_reviews;
-- +goose StatementEnd
