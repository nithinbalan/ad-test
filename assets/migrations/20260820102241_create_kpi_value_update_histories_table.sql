-- +goose Up
-- +goose StatementBegin
CREATE TABLE kpi_value_update_histories (
    id INT IDENTITY(1, 1) PRIMARY KEY,
    kpi_value_update_id INT NOT NULL,
    value DECIMAL(18, 2) NULL,
    unit_id INT NULL,
    remarks NVARCHAR(MAX) NULL,
    corrective_action NVARCHAR(MAX) NULL,
    created_at DATETIME NOT NULL,
    created_by_id INT NULL,
    CONSTRAINT FK_kpi_value_update_histories_kpi_value_updates
        FOREIGN KEY (kpi_value_update_id) REFERENCES kpi_value_updates(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_value_update_histories_units
        FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_value_update_histories_users
        FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.kpi_value_update_histories;
-- +goose StatementEnd
