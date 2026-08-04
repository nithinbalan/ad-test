-- +goose Up
-- +goose StatementBegin
CREATE TABLE review_types (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
        CONSTRAINT CHK_review_types_type
            CHECK (type IN ('user', 'privileged', 'service', 'external')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_review_types_status
            CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE review_types;
-- +goose StatementEnd
