-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_reviews (
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_register_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    reviewed_at DATETIME NOT NULL,
    reviewed_by INT NOT NULL,
    CONSTRAINT FK_compliance_reviews_compliance_register FOREIGN KEY (compliance_register_id) REFERENCES compliance_register (id),
    CONSTRAINT FK_compliance_reviews_users FOREIGN KEY (reviewed_by) REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_reviews;
-- +goose StatementEnd
