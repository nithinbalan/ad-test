-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_statuses (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
        CONSTRAINT CHK_audit_statuses_type CHECK (type IN ('pending', 'planned', 'in-progress', 'completed')),
    color_code VARCHAR(10) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_audit_statuses_status CHECK (status IN ('active', 'inactive'))
   
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_statuses;
-- +goose StatementEnd
