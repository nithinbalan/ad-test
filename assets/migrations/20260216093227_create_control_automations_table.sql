-- +goose Up
-- +goose StatementBegin
CREATE TABLE control_automations (
    id INT PRIMARY KEY IDENTITY(1,1),

    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,

    type NVARCHAR(50) NOT NULL,
    score DECIMAL(4,2) NOT NULL,

    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_automation_levels_status
        CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE control_automations;
-- +goose StatementEnd
