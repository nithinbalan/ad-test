-- +goose Up
-- +goose StatementBegin
ALTER TABLE project_decision_log_review_comments
    ADD comment_status_id INT NULL,
        level INT NULL,
        parent_comment_id INT NULL;

ALTER TABLE project_decision_log_review_comments
    ADD CONSTRAINT FK_project_decision_log_review_comments_statuses
        FOREIGN KEY (comment_status_id) REFERENCES comment_statuses(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE project_decision_log_review_comments
    ADD CONSTRAINT FK_project_decision_log_review_comments_parent
        FOREIGN KEY (parent_comment_id) REFERENCES project_decision_log_review_comments(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE project_decision_log_review_comments
    DROP CONSTRAINT FK_project_decision_log_review_comments_parent;

ALTER TABLE project_decision_log_review_comments
    DROP CONSTRAINT FK_project_decision_log_review_comments_statuses;

ALTER TABLE project_decision_log_review_comments
    DROP COLUMN comment_status_id,
                level,
                parent_comment_id;
-- +goose StatementEnd
