-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_post_review_users
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    incident_post_review_id INT NOT NULL,
    user_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_incident_post_review_users_incident_post_reviews FOREIGN KEY (incident_post_review_id) REFERENCES incident_post_reviews(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_post_review_users_users_1 FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_post_reviews_users_users_2 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,

);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_post_review_users;

-- +goose StatementEnd