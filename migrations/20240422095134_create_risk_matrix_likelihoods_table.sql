-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_matrix_likelihoods (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_matrix_id INT NOT NULL,
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    description_en NVARCHAR(max) NOT NULL,
    description_ar NVARCHAR(max) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_risk_matrix_likelihoods_status CHECK (status IN ('active', 'inactive')),
    CONSTRAINT FK_risk_matrix_likelihoods_risk_matrix FOREIGN KEY (risk_matrix_id) REFERENCES risk_matrix(id)
   
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_matrix_likelihoods;
-- +goose StatementEnd
