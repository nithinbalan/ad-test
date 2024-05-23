-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_post_reviews
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    incident_register_id INT NOT NULL,
    incident_post_review_effectiveness_id INT NOT NULL,
    post_review_date Date NOT NULL,
    summary NVARCHAR(MAX) NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_incident_post_reviews_incident_register FOREIGN KEY (incident_register_id) REFERENCES incident_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_post_reviews_incident_post_review_effectivenesses FOREIGN KEY (incident_post_review_effectiveness_id) REFERENCES incident_post_review_effectivenesses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_post_reviews_users_1 FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_post_reviews_users_2 FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_post_reviews;

-- +goose StatementEnd