-- +goose Up
-- +goose StatementBegin
CREATE TABLE kpi_value_update_history_inputs (
    id INT IDENTITY(1, 1) PRIMARY KEY,
    kpi_value_update_history_id INT NOT NULL,
    variable NVARCHAR(45) NOT NULL,
    input NVARCHAR(255) NOT NULL,
    value DECIMAL(18, 2) NOT NULL,
    created_by_id INT NULL,
    CONSTRAINT FK_kpi_value_update_history_inputs_histories
        FOREIGN KEY (kpi_value_update_history_id) REFERENCES kpi_value_update_histories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_value_update_history_inputs_users
        FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.kpi_value_update_history_inputs;
-- +goose StatementEnd
