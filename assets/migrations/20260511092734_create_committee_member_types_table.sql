-- +goose Up
-- +goose StatementBegin
CREATE TABLE committee_member_types (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en VARCHAR(100) NOT NULL,
    title_ar VARCHAR(100) NOT NULL,
    status VARCHAR(45) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_committee_member_types_status 
        CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.committee_member_types;
-- +goose StatementEnd