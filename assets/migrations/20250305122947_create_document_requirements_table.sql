-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_requirements (
    id INT PRIMARY KEY IDENTITY(1, 1),
    document_requirement_group_id INT,
    title_en NVARCHAR(250) NOT NULL,
    title_ar NVARCHAR(250) NOT NULL,  
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,  
    CONSTRAINT FK_document_requirements_requirement_groups FOREIGN KEY (document_requirement_group_id) REFERENCES document_requirement_groups(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_requirements;
-- +goose StatementEnd
