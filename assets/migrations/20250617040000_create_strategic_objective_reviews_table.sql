-- +goose Up
-- +goose StatementBegin

CREATE TABLE strategic_objective_reviews (
    id INT PRIMARY KEY IDENTITY(1,1),
    strategic_objective_id INT NOT NULL,
    review_date DATE NOT NULL,
    description NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_so_reviews_objectives FOREIGN KEY (strategic_objective_id) 
        REFERENCES strategic_objectives(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_so_reviews_users_created FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_so_reviews_users_updated FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE strategic_objective_reviews;

-- +goose StatementEnd
