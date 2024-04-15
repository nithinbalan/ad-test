-- +goose Up
-- +goose StatementBegin
CREATE TABLE strategic_objectives (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en VARCHAR(MAX) NOT NULL,
    title_ar NVARCHAR(MAX) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    created_by INT NOT NULL,
    updated_by INT,
    CONSTRAINT CHK_strategic_objectives_status CHECK (status IN ('active', 'inactive')),
    CONSTRAINT FK_strategic_objectives_created_by FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT FK_strategic_objectives_updated_by FOREIGN KEY (updated_by) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS strategic_objectives;
-- +goose StatementEnd
