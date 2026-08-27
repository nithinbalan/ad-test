-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_execution_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    type VARCHAR(40) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_penetration_test_execution_statuses_type CHECK (type IN ('planned', 'environment-preparation', 'in-progress', 'completed')),
    CONSTRAINT CHK_penetration_test_execution_statuses_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS penetration_test_execution_statuses;
-- +goose StatementEnd
