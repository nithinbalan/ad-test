-- +goose Up
-- +goose StatementBegin
CREATE TABLE control_reviews (
    id INT PRIMARY KEY IDENTITY(1, 1),
    control_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    reviewed_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    reviewed_by INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_control_reviews_controls FOREIGN KEY (control_id) REFERENCES controls (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_control_reviews_users FOREIGN KEY (reviewed_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE control_reviews;
-- +goose StatementEnd
