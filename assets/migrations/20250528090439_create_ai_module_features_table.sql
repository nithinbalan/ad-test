-- +goose Up
-- +goose StatementBegin
CREATE TABLE ai_module_features (
    id INT PRIMARY KEY IDENTITY(1,1),
    ai_module_id INT NOT NULL,
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NULL,
    description_en NVARCHAR(MAX),
    description_ar NVARCHAR(MAX),
    knowledge_center_resource_status_id INT,
    [order] INT DEFAULT 0,
    CONSTRAINT FK_ai_module_features_ai_module FOREIGN KEY (ai_module_id) REFERENCES ai_modules(id) ON DELETE CASCADE,
    CONSTRAINT FK_ai_module_features_knowledge_center_resource_status FOREIGN KEY (knowledge_center_resource_status_id) REFERENCES knowledge_center_resource_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ai_module_features;
-- +goose StatementEnd
