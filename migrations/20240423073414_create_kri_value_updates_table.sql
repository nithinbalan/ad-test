-- +goose Up
-- +goose StatementBegin
CREATE TABLE kri_value_updates (
    id INT PRIMARY KEY IDENTITY(1,1),
    kri_register_id INT NOT NULL,
    reviewed_date DATE NOT NULL,
    value VARCHAR(455),
    low_threshold int null,
    medium_threshold int null,
    high_threshold int null,
    remarks NVARCHAR(MAX),
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_kri_value_updates_kri_register FOREIGN KEY (kri_register_id) REFERENCES kri_register (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kri_value_updates_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kri_value_updates_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS kri_value_updates;
-- +goose StatementEnd
