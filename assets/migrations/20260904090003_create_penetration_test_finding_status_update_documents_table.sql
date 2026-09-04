-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_finding_status_update_documents (
    id INT PRIMARY KEY IDENTITY(1, 1),
    penetration_test_finding_status_update_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    CONSTRAINT FK_penetration_test_finding_status_update_documents_status_update
        FOREIGN KEY (penetration_test_finding_status_update_id) REFERENCES penetration_test_finding_status_updates(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_finding_status_update_documents_created_by
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_finding_status_update_documents_updated_by
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE penetration_test_finding_status_update_documents;
-- +goose StatementEnd
