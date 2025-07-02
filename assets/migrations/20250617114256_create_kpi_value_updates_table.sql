-- +goose Up
-- +goose StatementBegin
CREATE TABLE kpi_value_updates (
    id INT PRIMARY KEY IDENTITY(1,1),
    kpi_register_id INT NOT NULL,
    reviewed_date DATE NOT NULL,
    review_user_id INT NULL,
    value INT NULL,
    unit_id INT NULL,
    remarks NVARCHAR(MAX),
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_kpi_value_updates_kpi_register FOREIGN KEY (kpi_register_id) REFERENCES kpi_register (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_value_updates_users_3 FOREIGN KEY (review_user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_value_updates_units FOREIGN KEY (unit_id) REFERENCES units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_value_updates_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_value_updates_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS kpi_value_updates;
-- +goose StatementEnd
