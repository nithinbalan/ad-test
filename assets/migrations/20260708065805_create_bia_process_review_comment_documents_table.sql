-- +goose Up
-- +goose StatementBegin
CREATE TABLE bia_process_review_comment_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    bia_process_review_comment_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    CONSTRAINT FK_bia_process_rvcd_bia_process_review_comment FOREIGN KEY (bia_process_review_comment_id) REFERENCES bia_process_review_comments(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bia_process_rvcd_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bia_process_rvcd_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bia_process_review_comment_documents;
-- +goose StatementEnd
