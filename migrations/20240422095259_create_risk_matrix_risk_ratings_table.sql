-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_matrix_risk_ratings (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_matrix_id INT NOT NULL,
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(10) NOT NULL,
        CONSTRAINT CHK_risk_matrix_risk_ratings_type CHECK (type IN ('very_low','low', 'medium', 'high', 'very_high')),
    color_code VARCHAR(10) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT FK_risk_matrix_risk_ratings_risk_matrix FOREIGN KEY (risk_matrix_id) REFERENCES risk_matrix(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT CHK_risk_matrix_risk_ratings_status CHECK (status IN ('active', 'inactive'))
   
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_matrix_risk_ratings;
-- +goose StatementEnd
