-- +goose Up
-- +goose StatementBegin
CREATE TABLE knowledge_center_resource_statuses (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_knowledge_center_resource_statuses_status 
        CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE knowledge_center_resource_statuses;
-- +goose StatementEnd
