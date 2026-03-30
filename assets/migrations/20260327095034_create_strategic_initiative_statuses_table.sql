-- +goose Up
-- +goose StatementBegin
CREATE TABLE strategic_initiative_statuses (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL DEFAULT 'planned',
        CONSTRAINT CHK_strategic_initiative_statuses_type
            CHECK (type IN ('planned', 'in_progress', 'completed', 'on_hold')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_strategic_initiative_statuses_status
            CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.strategic_initiative_statuses;
-- +goose StatementEnd