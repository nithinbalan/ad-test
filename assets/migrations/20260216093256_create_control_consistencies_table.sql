-- +goose Up
-- +goose StatementBegin
CREATE TABLE control_consistencies (
    id INT PRIMARY KEY IDENTITY(1,1),

    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,

    type NVARCHAR(50) NOT NULL,
    score DECIMAL(4,2) NOT NULL,

    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_control_consistency_status
        CHECK (status IN ('active', 'inactive')),

    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    CONSTRAINT FK_control_consistencies_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE control_consistencies;
-- +goose StatementEnd
