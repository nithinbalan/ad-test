-- +goose Up
-- +goose StatementBegin
ALTER TABLE processes
ADD 
    risk_matrix_risk_rating_id INT,
    risk_matrix_risk_rating_value_id INT,
    CONSTRAINT FK_processes_risk_matrix_risk_ratings FOREIGN KEY (risk_matrix_risk_rating_id) REFERENCES risk_matrix_risk_ratings(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_processes_risk_matrix_risk_rating_values FOREIGN KEY (risk_matrix_risk_rating_value_id) REFERENCES risk_matrix_risk_rating_values(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE processes
DROP CONSTRAINT FK_processes_risk_matrix_risk_ratings,
DROP CONSTRAINT FK_processes_risk_matrix_risk_rating_values,
DROP COLUMN risk_matrix_risk_rating_id,
DROP COLUMN risk_matrix_risk_rating_value_id;
-- +goose StatementEnd
