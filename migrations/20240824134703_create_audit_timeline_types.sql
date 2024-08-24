-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_timeline_types (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL CHECK (type IN ('audited', 'fieldwork', 'draft_report', 'final_report', 'closure'))
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_timeline_types;
-- +goose StatementEnd
