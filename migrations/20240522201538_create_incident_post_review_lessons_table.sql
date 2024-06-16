-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_post_review_lessons
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    incident_post_review_id INT NOT NULL,
    lesson NVARCHAR(500) NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_incident_post_review_lessons_incident_post_reviews FOREIGN KEY (incident_post_review_id) REFERENCES incident_post_reviews(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_post_review_lessons_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION

);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.incident_post_review_lessons;
-- +goose StatementEnd
