-- +goose Up
-- +goose StatementBegin

CREATE TABLE bcm_process_types (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(MAX) NOT NULL,
    title_ar NVARCHAR(MAX) NOT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,   
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_bcm_process_types_status CHECK (status IN ('active', 'inactive')),

    CONSTRAINT FK_bcm_process_types_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcm_process_types_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE dbo.bcm_process_types;

-- +goose StatementEnd
