-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_change_request_initiate_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(10) NOT NULL,
        CONSTRAINT CHK_document_change_request_initiate_statuses_type CHECK (type IN ('initiated', 'approved', 'rejected')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_document_change_request_initiate_statuses_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_change_request_initiate_statuses;
-- +goose StatementEnd
