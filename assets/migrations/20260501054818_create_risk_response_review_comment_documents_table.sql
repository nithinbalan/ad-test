-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_response_review_comment_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_response_review_comment_id INT NOT NULL,
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
    CONSTRAINT FK_risk_response_review_comment_documents_risk_response_review_comment FOREIGN KEY (risk_response_review_comment_id) REFERENCES risk_response_review_comments(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_response_review_comment_documents_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_response_review_comment_documents_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_response_review_comment_documents;
-- +goose StatementEnd
