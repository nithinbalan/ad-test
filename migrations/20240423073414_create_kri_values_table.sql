-- +goose Up
-- +goose StatementBegin
CREATE TABLE kri_values (
    id INT PRIMARY KEY IDENTITY(1,1),
    kri_register_id INT NOT NULL,
    review_date DATE NOT NULL,
    value VARCHAR(455),
    low_threshold VARCHAR(455),
    medium_threshold VARCHAR(455),
    high_threshold VARCHAR(455),
    remarks NVARCHAR(MAX),
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    created_by INT NOT NULL,
    updated_by INT,
    CONSTRAINT FK_kri_values_kri_register FOREIGN KEY (kri_register_id) REFERENCES kri_register (id),
    CONSTRAINT FK_kri_values_users1 FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT FK_kri_values_users2 FOREIGN KEY (updated_by) REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS kri_values;
-- +goose StatementEnd
