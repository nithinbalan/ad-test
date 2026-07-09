-- +goose Up
-- +goose StatementBegin
CREATE TABLE vulnerability_remediation_users
(
    id INT PRIMARY KEY IDENTITY(1,1),

    vulnerability_id INT NOT NULL,
    user_id INT NOT NULL,

    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_vulnerability_remediation_users_vulnerability FOREIGN KEY (vulnerability_id)
        REFERENCES vulnerabilities(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_vulnerability_remediation_users_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_vulnerability_remediation_users_created_by FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE vulnerability_remediation_users;
-- +goose StatementEnd