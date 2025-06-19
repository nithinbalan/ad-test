-- +goose Up
-- +goose StatementBegin
CREATE TABLE business_continuity_plans (
    id INT PRIMARY KEY IDENTITY(1,1),
    reference_number VARCHAR(45) NOT NULL UNIQUE,
    business_function_id INT NULL,
    version VARCHAR(45) NULL,
    purpose NVARCHAR(MAX) NULL,
    scope NVARCHAR(MAX) NULL,
    objectives NVARCHAR(MAX) NULL,
    assumptions NVARCHAR(MAX) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,
    CONSTRAINT FK_business_continuity_plans_business_functions FOREIGN KEY (business_function_id) REFERENCES business_functions(id),
    CONSTRAINT FK_business_continuity_plans_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id),
    CONSTRAINT FK_business_continuity_plans_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE business_continuity_plans;
-- +goose StatementEnd
