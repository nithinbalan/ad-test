-- +goose Up
-- +goose StatementBegin
CREATE TABLE cs_log_soc_room_access_log_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    soc_room_access_log_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_cs_log_soc_room_access_log_documents_cs_log_soc_room_access_logs FOREIGN KEY (soc_room_access_log_id) REFERENCES cs_log_soc_room_access_logs(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_cs_log_soc_room_access_log_documents_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_cs_log_soc_room_access_log_documents_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.cs_log_soc_room_access_log_documents;
-- +goose StatementEnd
