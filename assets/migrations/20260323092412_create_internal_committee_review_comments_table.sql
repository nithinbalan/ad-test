-- +goose Up
-- +goose StatementBegin
CREATE TABLE internal_committee_review_comments
(
    id INT PRIMARY KEY IDENTITY(1,1),
    internal_committee_id INT NOT NULL,
    internal_committee_version_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    level INT NULL,
    comment_status_id INT NOT NULL,
    internal_committee_review_comment_id INT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_internal_committee_review_comments FOREIGN KEY (internal_committee_id) REFERENCES internal_committees(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_internal_committee_review_comments_version FOREIGN KEY (internal_committee_version_id) REFERENCES internal_committee_versions(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_review_comments_parent_comment FOREIGN KEY (internal_committee_review_comment_id) REFERENCES internal_committee_review_comments(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_review_comments_comment_statuses FOREIGN KEY (comment_status_id) REFERENCES comment_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_review_comments_users FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_review_comments_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE internal_committee_review_comments;
-- +goose StatementEnd
