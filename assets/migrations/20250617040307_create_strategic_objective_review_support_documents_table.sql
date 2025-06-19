-- +goose Up
-- +goose StatementBegin

CREATE TABLE strategic_objective_review_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    strategic_objective_review_id INT NOT NULL,
    url VARCHAR(1000) NOT NULL,
    name NVARCHAR(255) NOT NULL,
    size INT NOT NULL,
    token VARCHAR(255) NOT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_so_review_docs_reviews FOREIGN KEY (strategic_objective_review_id)
        REFERENCES strategic_objective_reviews(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE strategic_objective_review_documents;

-- +goose StatementEnd
