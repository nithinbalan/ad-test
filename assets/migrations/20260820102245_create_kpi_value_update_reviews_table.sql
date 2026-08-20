-- +goose Up
-- +goose StatementBegin
CREATE TABLE kpi_value_update_reviews (
    id INT IDENTITY(1, 1) PRIMARY KEY,
    kpi_value_update_id INT NOT NULL,
    kpi_review_status_id INT NOT NULL,
    comment NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_kpi_value_update_reviews_kpi_value_updates
        FOREIGN KEY (kpi_value_update_id) REFERENCES kpi_value_updates(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_value_update_reviews_kpi_review_statuses
        FOREIGN KEY (kpi_review_status_id) REFERENCES kpi_review_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_value_update_reviews_users
        FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.kpi_value_update_reviews;
-- +goose StatementEnd
