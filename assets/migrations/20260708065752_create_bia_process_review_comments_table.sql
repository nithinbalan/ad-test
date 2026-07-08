-- +goose Up
-- +goose StatementBegin
CREATE TABLE bia_process_review_comments
(
    id INT PRIMARY KEY IDENTITY(1,1),
    bia_process_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    level INT NULL,
    comment_status_id INT NOT NULL,
    parent_comment_id INT NULL,
    is_edited BIT NOT NULL DEFAULT 0,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_bia_process_rvc_bia_process
        FOREIGN KEY (bia_process_id)
        REFERENCES bia_processes(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT FK_bia_process_rvc_parent
        FOREIGN KEY (parent_comment_id)
        REFERENCES bia_process_review_comments(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_bia_process_rvc_comment_statuses
        FOREIGN KEY (comment_status_id)
        REFERENCES comment_statuses(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_bia_process_rvc_users_1
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_bia_process_rvc_users_2
        FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE bia_process_review_comments;
-- +goose StatementEnd
