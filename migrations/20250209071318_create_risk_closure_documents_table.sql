-- +goose Up
-- +goose StatementBegin

CREATE TABLE risk_closure_documents (
    id INT PRIMARY KEY IDENTITY(1, 1),
    risk_closure_id INT NOT NULL,
    url NVARCHAR(1000) NOT NULL,           
    name NVARCHAR(255) NOT NULL,           
    size INT NOT NULL,         
    token VARCHAR(255) NOT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_risk_closure_documents_risk_closures FOREIGN KEY (risk_closure_id) REFERENCES risk_closures (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_risk_closure_documents_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id),
    CONSTRAINT FK_risk_closure_documents_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS risk_closure_documents;

-- +goose StatementEnd
