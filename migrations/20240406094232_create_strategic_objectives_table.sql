-- +goose Up
-- +goose StatementBegin
CREATE TABLE strategic_objectives (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(MAX) NOT NULL,
    title_ar NVARCHAR(MAX) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_strategic_objectives_status CHECK (status IN ('active', 'inactive')),
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    created_by INT NOT NULL,
    updated_by INT,
    CONSTRAINT FK_strategic_objectives_users_1 FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategic_objectives_users_2 FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS strategic_objectives;
-- +goose StatementEnd
