-- +goose Up
-- +goose StatementBegin
CREATE TABLE kpi_value_update_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    kpi_value_update_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name VARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    type VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    CONSTRAINT FK_kpi_value_update_documents FOREIGN KEY (kpi_value_update_id) REFERENCES kpi_value_updates(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE kpi_value_update_documents;
-- +goose StatementEnd
