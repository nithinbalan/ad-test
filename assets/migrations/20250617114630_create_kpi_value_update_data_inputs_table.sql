-- +goose Up
-- +goose StatementBegin
CREATE TABLE kpi_value_update_data_inputs (
    id INT PRIMARY KEY IDENTITY(1,1),
    kpi_value_update_id INT NOT NULL,
    variable NVARCHAR(45) NOT NULL,
    input NVARCHAR(255) NOT NULL,
    value INT NOT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NULL,

    -- Foreign Key Constraints
    CONSTRAINT FK_kpi_value_update_data_inputs_kpi_value_update FOREIGN KEY (kpi_value_update_id) 
        REFERENCES kpi_value_updates(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_value_update_data_inputs_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_value_update_data_inputs_users_2 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE kpi_value_update_data_inputs;
-- +goose StatementEnd
