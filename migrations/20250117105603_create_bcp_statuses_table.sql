-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcp_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL DEFAULT 'draft',
        CONSTRAINT CHK_bcp_statuses_type CHECK (type IN ('draft', 'in-review', 'published')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_bcp_statuses_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.bcp_statuses;
-- +goose StatementEnd

