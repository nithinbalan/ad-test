-- +goose Up
-- +goose StatementBegin
CREATE TABLE kri_value_update_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    kri_value_update_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name VARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    type VARCHAR(255) DEFAULT NULL,
    icon VARCHAR(10) NOT NULL, -- Assuming icon names like 'pdf', 'excel', etc.
        CONSTRAINT CHK_kri_value_update_documents_icon CHECK (icon IN ('pdf', 'excel', 'ppt', 'image', 'text', 'other')),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    CONSTRAINT FK_kri_value_update_documents FOREIGN KEY (kri_value_update_id) REFERENCES kri_value_updates(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE kri_value_update_documents;
-- +goose StatementEnd
