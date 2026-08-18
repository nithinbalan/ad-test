-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_closure_review_comments
(
    id INT PRIMARY KEY IDENTITY(1,1),
    incident_closure_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    level INT NULL,
    comment_status_id INT NOT NULL,
    is_edited BIT NOT NULL DEFAULT 0,
    incident_closure_review_comment_id INT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_incident_closure_review_comments FOREIGN KEY (incident_closure_id) REFERENCES incident_closures(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_incident_closure_review_comments_parent_comment FOREIGN KEY (incident_closure_review_comment_id) REFERENCES incident_closure_review_comments(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_closure_review_comments_comment_statuses FOREIGN KEY (comment_status_id) REFERENCES comment_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_closure_review_comments_users FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_closure_review_comments_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_closure_review_comments;
-- +goose StatementEnd
