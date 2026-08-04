-- +goose Up
-- +goose StatementBegin
CREATE TABLE identity_review_users (
    id INT PRIMARY KEY IDENTITY(1, 1),
    username NVARCHAR(255) NOT NULL,
    department_id INT NOT NULL,
    privilege_id INT NULL,
    account_type_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,

	CONSTRAINT FK_identity_review_users_departments
		FOREIGN KEY (department_id) REFERENCES organization_units(id)
		ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_users_privileged_accounts
        FOREIGN KEY (privilege_id) REFERENCES privileged_accounts(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_users_account_types
        FOREIGN KEY (account_type_id) REFERENCES account_types(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_users_1
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_users_2
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE  dbo.identity_review_users;
-- +goose StatementEnd

