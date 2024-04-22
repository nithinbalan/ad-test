-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_matrix_risk_rating_values (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_rating_id INT NOT NULL,
    value INT NOT NULL,
    CONSTRAINT FK_risk_matrix_risk_rating_values_risk_matrix_risk_ratings FOREIGN KEY (risk_rating_id) REFERENCES risk_matrix_risk_ratings(id)
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_matrix_risk_rating_values;
-- +goose StatementEnd
