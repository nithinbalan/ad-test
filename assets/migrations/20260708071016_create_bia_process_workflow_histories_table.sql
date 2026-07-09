-- +goose Up
-- +goose StatementBegin
CREATE TABLE bia_process_workflow_histories (
    id INT PRIMARY KEY IDENTITY(1,1),
    bia_process_id INT NOT NULL,
    level INT NOT NULL,
    workflow_status_id INT NOT NULL,
    comment NVARCHAR(MAX),
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,

    CONSTRAINT FK_bia_process_workflow_histories_workflow
        FOREIGN KEY (bia_process_id)
        REFERENCES bia_processes(id)
       ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_bia_process_workflow_histories_status
        FOREIGN KEY (workflow_status_id)
        REFERENCES workflow_statuses(id)
       ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_bia_process_workflow_histories_users_1
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
       ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_bia_process_workflow_histories_users_2
        FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
       ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bia_process_workflow_histories;
-- +goose StatementEnd