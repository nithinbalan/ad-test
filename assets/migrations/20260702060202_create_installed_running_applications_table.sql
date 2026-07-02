-- +goose Up
-- +goose StatementBegin
CREATE TABLE installed_running_applications(
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,

    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_installed_running_applications_status
        CHECK (status IN ('active', 'inactive'))
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE installed_running_applications;
-- +goose StatementEnd
