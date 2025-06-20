-- +goose Up
-- +goose StatementBegin
CREATE TABLE knowledge_center_resources (
    id INT PRIMARY KEY IDENTITY(1,1),
    module_id INT NOT NULL,
    knowledge_center_resource_status_id INT NOT NULL,
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    icon NVARCHAR(255),
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_knowledge_center_resources_modules
        FOREIGN KEY (module_id) REFERENCES modules(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_knowledge_center_resources_statuses
        FOREIGN KEY (knowledge_center_resource_status_id) REFERENCES knowledge_center_resource_statuses(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE knowledge_center_resources;
-- +goose StatementEnd
