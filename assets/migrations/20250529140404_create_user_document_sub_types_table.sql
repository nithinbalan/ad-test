-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_document_sub_types (
    id INT PRIMARY KEY IDENTITY(1, 1),
    user_document_type_id INT NULL,
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_user_document_sub_types_status CHECK (status IN ('active', 'inactive')),
    CONSTRAINT FK_user_document_sub_types_user_document_types FOREIGN KEY (user_document_type_id) REFERENCES user_document_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.user_document_sub_types;
-- +goose StatementEnd
